"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || (function () {
    var ownKeys = function(o) {
        ownKeys = Object.getOwnPropertyNames || function (o) {
            var ar = [];
            for (var k in o) if (Object.prototype.hasOwnProperty.call(o, k)) ar[ar.length] = k;
            return ar;
        };
        return ownKeys(o);
    };
    return function (mod) {
        if (mod && mod.__esModule) return mod;
        var result = {};
        if (mod != null) for (var k = ownKeys(mod), i = 0; i < k.length; i++) if (k[i] !== "default") __createBinding(result, mod, k[i]);
        __setModuleDefault(result, mod);
        return result;
    };
})();
Object.defineProperty(exports, "__esModule", { value: true });
exports.activate = activate;
exports.deactivate = deactivate;
const vscode = __importStar(require("vscode"));
const cp = __importStar(require("child_process"));
const net = __importStar(require("net"));
const rpc_1 = require("./rpc");
let client;
let out;
async function ensureClient(context) {
    if (client)
        return client;
    const cfg = vscode.workspace.getConfiguration('eve');
    const cliPath = (0, rpc_1.resolveCliPath)(cfg);
    const mode = cfg.get('serverMode') || 'stdio';
    out.appendLine(`[EVE] starting client mode=${mode}`);
    if (mode === 'http') {
        let port = cfg.get('httpPort') || 0;
        if (!port) {
            const srv = net.createServer();
            await new Promise(res => srv.listen(0, res));
            port = srv.address().port;
            srv.close();
        }
        const proc = cp.spawn(cliPath, ['--serve', `--port=${port}`, ...((cfg.get('args')) || [])], { stdio: 'inherit' });
        proc.on('close', code => out.appendLine(`[EVE] http server exited: ${code}`));
        client = new rpc_1.HttpRpcClient('127.0.0.1', port);
        return client;
    }
    client = new rpc_1.StdioRpcClient(cliPath, ['--daemon', ...((cfg.get('args')) || [])], out);
    return client;
}
async function chatCmd() {
    const c = await ensureClient();
    const prompt = await vscode.window.showInputBox({ prompt: 'EVE Prompt' });
    if (!prompt)
        return;
    const doc = vscode.window.activeTextEditor?.document;
    const context = doc ? { uri: doc.uri.toString(), text: doc.getText() } : {};
    const resp = await c.request('chat', { prompt, context });
    out.appendLine(`\n[EVE Chat]\n${resp}\n`);
}
async function planCmd() {
    const c = await ensureClient();
    const targets = vscode.workspace.workspaceFolders?.map(f => f.uri.fsPath) || [];
    const plan = await c.request('planEdits', { targets });
    out.appendLine(`[EVE Plan] ${JSON.stringify(plan, null, 2)}`);
}
async function applyCmd() {
    const c = await ensureClient();
    const plan = await c.request('planEdits', {});
    const edits = plan?.edits || [];
    const we = new vscode.WorkspaceEdit();
    for (const e of edits) {
        const uri = vscode.Uri.file(e.path);
        if (e.start && e.end) {
            we.replace(uri, new vscode.Range(new vscode.Position(e.start.line, e.start.character), new vscode.Position(e.end.line, e.end.character)), e.newText);
        }
        else {
            const doc = await vscode.workspace.openTextDocument(uri);
            we.replace(uri, new vscode.Range(new vscode.Position(0, 0), doc.lineAt(doc.lineCount - 1).range.end), e.newText);
        }
    }
    const ok = await vscode.workspace.applyEdit(we);
    out.appendLine(`[EVE Apply] ${ok ? 'Applied' : 'Failed'}`);
}
async function openWebCmd(context) {
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
            const resp = await c.request('chat', { prompt: msg.prompt });
            panel.webview.postMessage({ type: 'chatResult', text: resp });
        }
        else if (msg.type === 'plan') {
            const plan = await c.request('planEdits', {});
            panel.webview.postMessage({ type: 'planResult', plan });
        }
        else if (msg.type === 'apply') {
            await applyCmd();
            panel.webview.postMessage({ type: 'applyDone' });
        }
    });
}
function activate(context) {
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
function deactivate() { client?.dispose(); }
//# sourceMappingURL=extension.js.map