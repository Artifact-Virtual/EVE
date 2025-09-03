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
exports.HttpRpcClient = exports.StdioRpcClient = void 0;
exports.resolveCliPath = resolveCliPath;
const cp = __importStar(require("child_process"));
const http = __importStar(require("http"));
const crypto = __importStar(require("crypto"));
class StdioRpcClient {
    constructor(cmd, args, out) {
        this.out = out;
        this.pending = new Map();
        this.buf = '';
        this.proc = cp.spawn(cmd, args, { stdio: 'pipe' });
        this.proc.stdout.setEncoding('utf8');
        this.proc.stdout.on('data', (chunk) => {
            this.buf += chunk;
            let idx;
            while ((idx = this.buf.indexOf('\n')) >= 0) {
                const line = this.buf.slice(0, idx);
                this.buf = this.buf.slice(idx + 1);
                try {
                    const msg = JSON.parse(line);
                    if (msg.id && this.pending.has(msg.id)) {
                        const p = this.pending.get(msg.id);
                        this.pending.delete(msg.id);
                        if ('error' in msg)
                            p.reject(new Error(msg.error?.message ?? 'Unknown error'));
                        else
                            p.resolve(msg.result);
                    }
                    else {
                        this.out.appendLine(`[EVE unsolicited] ${line}`);
                    }
                }
                catch (e) {
                    this.out.appendLine(`[EVE parse error] ${String(e)} :: ${line}`);
                }
            }
        });
        this.proc.stderr.setEncoding('utf8');
        this.proc.stderr.on('data', d => this.out.append(d.toString()));
        this.proc.on('close', c => this.out.appendLine(`EVE exited: ${c}`));
    }
    request(method, params) {
        const id = crypto.randomUUID();
        const payload = JSON.stringify({ jsonrpc: '2.0', id, method, params }) + '\n';
        return new Promise((resolve, reject) => {
            this.pending.set(id, { resolve, reject });
            this.proc.stdin.write(payload, 'utf8', (err) => {
                if (err) {
                    this.pending.delete(id);
                    reject(err);
                }
            });
        });
    }
    dispose() {
        try {
            this.proc.kill();
        }
        catch { }
        this.pending.forEach(p => p.reject(new Error('Disposed')));
        this.pending.clear();
    }
}
exports.StdioRpcClient = StdioRpcClient;
class HttpRpcClient {
    constructor(host, port) {
        this.host = host;
        this.port = port;
    }
    request(method, params) {
        const body = JSON.stringify({ method, params });
        return new Promise((resolve, reject) => {
            const req = http.request({ host: this.host, port: this.port, path: '/rpc', method: 'POST',
                headers: { 'content-type': 'application/json', 'content-length': Buffer.byteLength(body) } }, (res) => {
                const chunks = [];
                res.on('data', d => chunks.push(d));
                res.on('end', () => {
                    try {
                        resolve(JSON.parse(Buffer.concat(chunks).toString('utf8')));
                    }
                    catch (e) {
                        reject(e);
                    }
                });
            });
            req.on('error', reject);
            req.write(body);
            req.end();
        });
    }
    dispose() { }
}
exports.HttpRpcClient = HttpRpcClient;
function resolveCliPath(cfg) {
    let p = cfg.get('cliPath') || 'eve';
    if (process.platform === 'win32' && !p.endsWith('.exe'))
        p += '.exe';
    return p;
}
//# sourceMappingURL=rpc.js.map