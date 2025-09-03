import * as vscode from 'vscode';
import * as cp from 'child_process';
import * as net from 'net';
import { HttpRpcClient, RpcClient, StdioRpcClient, resolveCliPath } from './rpc';

let client: RpcClient | undefined;
let out: vscode.OutputChannel;

async function ensureClient(context?: vscode.ExtensionContext): Promise<RpcClient> {
  if (client) return client;
  const cfg = vscode.workspace.getConfiguration('eve');
  const cliPath = resolveCliPath(cfg);
  const mode = cfg.get<string>('serverMode') || 'stdio';

  out.appendLine(`[EVE] starting client mode=${mode}`);

  if (mode === 'http') {
    let port = cfg.get<number>('httpPort') || 0;
    if (!port) {
      const srv = net.createServer();
      await new Promise<void>(res => srv.listen(0, res));
      port = (srv.address() as any).port; srv.close();
    }
    const proc = cp.spawn(cliPath, ['--serve', `--port=${port}`, ...((cfg.get<string[]>('args'))||[])], {stdio: 'inherit'});
    proc.on('close', code => out.appendLine(`[EVE] http server exited: ${code}`));
    client = new HttpRpcClient('127.0.0.1', port);
    return client;
  }

  client = new StdioRpcClient(cliPath, ['--daemon', ...((cfg.get<string[]>('args'))||[])], out);
  return client;
}

async function chatCmd() {
  const c = await ensureClient();
  const prompt = await vscode.window.showInputBox({prompt: 'EVE Prompt'});
  if (!prompt) return;
  const doc = vscode.window.activeTextEditor?.document;
  const context = doc ? {uri: doc.uri.toString(), text: doc.getText()} : {};
  const resp = await c.request<string>('chat', {prompt, context});
  out.appendLine(`\n[EVE Chat]\n${resp}\n`);
}

async function planCmd() {
  const c = await ensureClient();
  const targets = vscode.workspace.workspaceFolders?.map(f => f.uri.fsPath) || [];
  const plan = await c.request<any>('planEdits', {targets});
  out.appendLine(`[EVE Plan] ${JSON.stringify(plan, null, 2)}`);
}

async function applyCmd() {
  const c = await ensureClient();
  const plan = await c.request<any>('planEdits', {});
  const edits: any[] = plan?.edits || [];
  const we = new vscode.WorkspaceEdit();
  for (const e of edits) {
    const uri = vscode.Uri.file(e.path);
    if (e.start && e.end) {
      we.replace(uri, new vscode.Range(
        new vscode.Position(e.start.line, e.start.character),
        new vscode.Position(e.end.line, e.end.character)
      ), e.newText);
    } else {
      const doc = await vscode.workspace.openTextDocument(uri);
      we.replace(uri, new vscode.Range(new vscode.Position(0,0), doc.lineAt(doc.lineCount-1).range.end), e.newText);
    }
  }
  const ok = await vscode.workspace.applyEdit(we);
  out.appendLine(`[EVE Apply] ${ok ? 'Applied' : 'Failed'}`);
}

async function openWebCmd(context: vscode.ExtensionContext) {
  const panel = vscode.window.createWebviewPanel('eveWeb', 'EVE WebUI', vscode.ViewColumn.One, {
    enableScripts: true,
    localResourceRoots: [vscode.Uri.joinPath(context.extensionUri, 'media')]
  });
  const scriptUri = panel.webview.asWebviewUri(vscode.Uri.joinPath(context.extensionUri, 'media', 'main.js'));
  panel.webview.html =
    `<!doctype html><html><head><meta charset='utf-8'>
      <meta http-equiv='Content-Security-Policy' content="default-src 'none'; img-src data:; script-src 'nonce-xyz'; style-src 'unsafe-inline';">
     </head><body><div id='app'></div><script nonce='xyz' src='${scriptUri}'></script></body></html>`;

  const c = await ensureClient(context);
  panel.webview.onDidReceiveMessage(async (msg) => {
    if (msg.type === 'chat') {
      const resp = await c.request<string>('chat', {prompt: msg.prompt});
      panel.webview.postMessage({type:'chatResult', text: resp});
    } else if (msg.type === 'plan') {
      const plan = await c.request<any>('planEdits', {});
      panel.webview.postMessage({type:'planResult', plan});
    } else if (msg.type === 'apply') {
      await applyCmd();
      panel.webview.postMessage({type:'applyDone'});
    }
  });
}

export function activate(context: vscode.ExtensionContext) {
  out = vscode.window.createOutputChannel('EVE');
  context.subscriptions.push(out);
  context.subscriptions.push(vscode.commands.registerCommand('eve.chat', chatCmd));
  context.subscriptions.push(vscode.commands.registerCommand('eve.plan', planCmd));
  context.subscriptions.push(vscode.commands.registerCommand('eve.apply', () => applyCmd()));
  context.subscriptions.push(vscode.commands.registerCommand('eve.openWeb', () => openWebCmd(context)));
  context.subscriptions.push(vscode.commands.registerCommand('eve.startDaemon', async () => { await ensureClient(context); out.show(); }));
  context.subscriptions.push(vscode.commands.registerCommand('eve.stopDaemon', async () => { client?.dispose(); client = undefined; out.appendLine('[EVE] Daemon stopped'); }));
  out.appendLine('EVE extension activated');
}

export function deactivate() { client?.dispose(); }
