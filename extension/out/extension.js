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
const rpc_1 = require("./rpc");
let out;
async function chatCmd() {
    const cfg = vscode.workspace.getConfiguration('eve');
    const cliPath = (0, rpc_1.resolveCliPath)(cfg);
    const prompt = await vscode.window.showInputBox({ prompt: 'EVE Prompt' });
    if (!prompt)
        return;
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
        const proc = cp.spawn(cliPath, [...((cfg.get('args')) || [])], {
            stdio: 'pipe',
            cwd: vscode.workspace.workspaceFolders?.[0]?.uri.fsPath
        });
        proc.stdin.write(contextualPrompt + '\n');
        proc.stdin.end();
        let output = '';
        proc.stdout.setEncoding('utf8');
        proc.stdout.on('data', (data) => {
            output += data;
            out.append(data);
        });
        proc.stderr.setEncoding('utf8');
        proc.stderr.on('data', (data) => {
            out.append(`[EVE Error] ${data}`);
        });
        proc.on('close', (code) => {
            out.appendLine(`\n[EVE] Process exited with code ${code}`);
        });
        proc.on('error', (err) => {
            out.appendLine(`[EVE] Error: ${err.message}`);
            vscode.window.showErrorMessage(`EVE Error: ${err.message}`);
        });
    }
    catch (error) {
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
    const cliPath = (0, rpc_1.resolveCliPath)(cfg);
    try {
        out.appendLine(`[EVE] Quick fix analysis for ${document.fileName}`);
        out.show();
        const proc = cp.spawn(cliPath, [...((cfg.get('args')) || [])], {
            stdio: 'pipe',
            cwd: vscode.workspace.workspaceFolders?.[0]?.uri.fsPath
        });
        proc.stdin.write(prompt + '\n');
        proc.stdin.end();
        proc.stdout.setEncoding('utf8');
        proc.stdout.on('data', (data) => {
            out.append(data);
        });
        proc.stderr.setEncoding('utf8');
        proc.stderr.on('data', (data) => {
            out.append(`[EVE Error] ${data}`);
        });
        proc.on('close', (code) => {
            out.appendLine(`\n[EVE] Analysis complete (exit code ${code})`);
        });
    }
    catch (error) {
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
    const cliPath = (0, rpc_1.resolveCliPath)(cfg);
    try {
        out.appendLine(`[EVE] Explaining code from ${document.fileName}`);
        out.show();
        const proc = cp.spawn(cliPath, [...((cfg.get('args')) || [])], {
            stdio: 'pipe',
            cwd: vscode.workspace.workspaceFolders?.[0]?.uri.fsPath
        });
        proc.stdin.write(prompt + '\n');
        proc.stdin.end();
        proc.stdout.setEncoding('utf8');
        proc.stdout.on('data', (data) => {
            out.append(data);
        });
        proc.stderr.setEncoding('utf8');
        proc.stderr.on('data', (data) => {
            out.append(`[EVE Error] ${data}`);
        });
        proc.on('close', (code) => {
            out.appendLine(`\n[EVE] Explanation complete (exit code ${code})`);
        });
    }
    catch (error) {
        const errorMsg = error instanceof Error ? error.message : String(error);
        out.appendLine(`[EVE] Error: ${errorMsg}`);
        vscode.window.showErrorMessage(`EVE Error: ${errorMsg}`);
    }
}
function activate(context) {
    out = vscode.window.createOutputChannel('EVE');
    context.subscriptions.push(out);
    context.subscriptions.push(vscode.commands.registerCommand('eve.chat', chatCmd));
    context.subscriptions.push(vscode.commands.registerCommand('eve.quickFix', quickFixCmd));
    context.subscriptions.push(vscode.commands.registerCommand('eve.explainCode', explainCodeCmd));
    out.appendLine('EVE extension activated');
    out.appendLine('Available commands: EVE: Chat, EVE: Quick Fix, EVE: Explain Code');
}
function deactivate() {
    if (out) {
        out.appendLine('EVE extension deactivated');
    }
}
//# sourceMappingURL=extension.js.map