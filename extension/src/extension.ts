import * as vscode from 'vscode';
import * as cp from 'child_process';
import { resolveCliPath } from './rpc';

let out: vscode.OutputChannel;

async function chatCmd() {
  const cfg = vscode.workspace.getConfiguration('eve');
  const cliPath = resolveCliPath(cfg);
  
  const prompt = await vscode.window.showInputBox({prompt: 'EVE Prompt'});
  if (!prompt) return;
  
  // Get active document context if available
  const doc = vscode.window.activeTextEditor?.document;
  let contextualPrompt = prompt;
  if (doc) {
    const fileName = doc.fileName;
    const selectedText = vscode.window.activeTextEditor?.selection && !vscode.window.activeTextEditor.selection.isEmpty 
      ? doc.getText(vscode.window.activeTextEditor.selection)
      : doc.getText();
    
    contextualPrompt = `File: ${fileName}\n\nContent:\n${selectedText}\n\nUser request: ${prompt}`;
  }
  
  try {
    out.appendLine(`[EVE] Starting chat with prompt: ${prompt}`);
    out.show();
    
    const proc = cp.spawn(cliPath, [...((cfg.get<string[]>('args'))||[])], {
      stdio: 'pipe',
      cwd: vscode.workspace.workspaceFolders?.[0]?.uri.fsPath
    });
    
    proc.stdin.write(contextualPrompt + '\n');
    proc.stdin.end();
    
    let output = '';
    proc.stdout.setEncoding('utf8');
    proc.stdout.on('data', (data: string) => {
      output += data;
      out.append(data);
    });
    
    proc.stderr.setEncoding('utf8');
    proc.stderr.on('data', (data: string) => {
      out.append(`[EVE Error] ${data}`);
    });
    
    proc.on('close', (code) => {
      out.appendLine(`\n[EVE] Process exited with code ${code}`);
    });
    
    proc.on('error', (err) => {
      out.appendLine(`[EVE] Error: ${err.message}`);
      vscode.window.showErrorMessage(`EVE Error: ${err.message}`);
    });
    
  } catch (error) {
    const errorMsg = error instanceof Error ? error.message : String(error);
    out.appendLine(`[EVE] Error: ${errorMsg}`);
    vscode.window.showErrorMessage(`EVE Error: ${errorMsg}`);
  }
}

async function quickFixCmd() {
  const editor = vscode.window.activeTextEditor;
  if (!editor) {
    vscode.window.showWarningMessage('No active editor found');
    return;
  }
  
  const document = editor.document;
  const selection = editor.selection;
  const selectedText = selection.isEmpty ? document.getText() : document.getText(selection);
  
  const prompt = `Please analyze and fix any issues in this code:\n\n${selectedText}`;
  
  // Simulate chat command with auto-generated prompt
  const cfg = vscode.workspace.getConfiguration('eve');
  const cliPath = resolveCliPath(cfg);
  
  try {
    out.appendLine(`[EVE] Quick fix analysis for ${document.fileName}`);
    out.show();
    
    const proc = cp.spawn(cliPath, [...((cfg.get<string[]>('args'))||[])], {
      stdio: 'pipe',
      cwd: vscode.workspace.workspaceFolders?.[0]?.uri.fsPath
    });
    
    proc.stdin.write(prompt + '\n');
    proc.stdin.end();
    
    proc.stdout.setEncoding('utf8');
    proc.stdout.on('data', (data: string) => {
      out.append(data);
    });
    
    proc.stderr.setEncoding('utf8');
    proc.stderr.on('data', (data: string) => {
      out.append(`[EVE Error] ${data}`);
    });
    
    proc.on('close', (code) => {
      out.appendLine(`\n[EVE] Analysis complete (exit code ${code})`);
    });
    
  } catch (error) {
    const errorMsg = error instanceof Error ? error.message : String(error);
    out.appendLine(`[EVE] Error: ${errorMsg}`);
    vscode.window.showErrorMessage(`EVE Error: ${errorMsg}`);
  }
}

async function explainCodeCmd() {
  const editor = vscode.window.activeTextEditor;
  if (!editor) {
    vscode.window.showWarningMessage('No active editor found');
    return;
  }
  
  const document = editor.document;
  const selection = editor.selection;
  const selectedText = selection.isEmpty ? document.getText() : document.getText(selection);
  
  const prompt = `Please explain what this code does:\n\n${selectedText}`;
  
  const cfg = vscode.workspace.getConfiguration('eve');
  const cliPath = resolveCliPath(cfg);
  
  try {
    out.appendLine(`[EVE] Explaining code from ${document.fileName}`);
    out.show();
    
    const proc = cp.spawn(cliPath, [...((cfg.get<string[]>('args'))||[])], {
      stdio: 'pipe',
      cwd: vscode.workspace.workspaceFolders?.[0]?.uri.fsPath
    });
    
    proc.stdin.write(prompt + '\n');
    proc.stdin.end();
    
    proc.stdout.setEncoding('utf8');
    proc.stdout.on('data', (data: string) => {
      out.append(data);
    });
    
    proc.stderr.setEncoding('utf8');
    proc.stderr.on('data', (data: string) => {
      out.append(`[EVE Error] ${data}`);
    });
    
    proc.on('close', (code) => {
      out.appendLine(`\n[EVE] Explanation complete (exit code ${code})`);
    });
    
  } catch (error) {
    const errorMsg = error instanceof Error ? error.message : String(error);
    out.appendLine(`[EVE] Error: ${errorMsg}`);
    vscode.window.showErrorMessage(`EVE Error: ${errorMsg}`);
  }
}

export function activate(context: vscode.ExtensionContext) {
  out = vscode.window.createOutputChannel('EVE');
  context.subscriptions.push(out);
  
  context.subscriptions.push(vscode.commands.registerCommand('eve.chat', chatCmd));
  context.subscriptions.push(vscode.commands.registerCommand('eve.quickFix', quickFixCmd));
  context.subscriptions.push(vscode.commands.registerCommand('eve.explainCode', explainCodeCmd));
  
  out.appendLine('EVE extension activated');
  out.appendLine('Available commands: EVE: Chat, EVE: Quick Fix, EVE: Explain Code');
}

export function deactivate() {
  if (out) {
    out.appendLine('EVE extension deactivated');
  }
}
