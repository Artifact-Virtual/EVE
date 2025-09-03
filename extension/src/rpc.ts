import * as cp from 'child_process';
import * as http from 'http';
import * as crypto from 'crypto';
import * as vscode from 'vscode';

type Json = any;

export interface RpcClient {
  request<T=any>(method: string, params?: Json): Promise<T>;
  dispose(): void;
}

export class StdioRpcClient implements RpcClient {
  private proc: cp.ChildProcessWithoutNullStreams;
  private pending = new Map<string, {resolve:(v:any)=>void, reject:(e:any)=>void}>();
  private buf = '';

  constructor(cmd: string, args: string[], private out: vscode.OutputChannel) {
    this.proc = cp.spawn(cmd, args, {stdio: 'pipe'});
    this.proc.stdout.setEncoding('utf8');
    this.proc.stdout.on('data', (chunk: string) => {
      this.buf += chunk;
      let idx: number;
      while ((idx = this.buf.indexOf('\n')) >= 0) {
        const line = this.buf.slice(0, idx);
        this.buf = this.buf.slice(idx + 1);
        try {
          const msg = JSON.parse(line);
          if (msg.id && this.pending.has(msg.id)) {
            const p = this.pending.get(msg.id)!;
            this.pending.delete(msg.id);
            if ('error' in msg) p.reject(new Error(msg.error?.message ?? 'Unknown error'));
            else p.resolve(msg.result);
          } else {
            this.out.appendLine(`[EVE unsolicited] ${line}`);
          }
        } catch (e) {
          this.out.appendLine(`[EVE parse error] ${String(e)} :: ${line}`);
        }
      }
    });
    this.proc.stderr.setEncoding('utf8');
    this.proc.stderr.on('data', d => this.out.append(d.toString()));
    this.proc.on('close', c => this.out.appendLine(`EVE exited: ${c}`));
  }

  request<T=any>(method: string, params?: Json): Promise<T> {
    const id = crypto.randomUUID();
    const payload = JSON.stringify({jsonrpc:'2.0', id, method, params}) + '\n';
    return new Promise<T>((resolve, reject) => {
      this.pending.set(id, {resolve, reject});
      this.proc.stdin.write(payload, 'utf8', (err) => {
        if (err) {
          this.pending.delete(id);
          reject(err);
        }
      });
    });
  }

  dispose(): void {
    try { this.proc.kill(); } catch {}
    this.pending.forEach(p => p.reject(new Error('Disposed')));
    this.pending.clear();
  }
}

export class HttpRpcClient implements RpcClient {
  constructor(private host: string, private port: number) {}
  request<T=any>(method: string, params?: Json): Promise<T> {
    const body = JSON.stringify({method, params});
    return new Promise<T>((resolve, reject) => {
      const req = http.request(
        {host:this.host, port:this.port, path:'/rpc', method:'POST',
         headers:{'content-type':'application/json','content-length':Buffer.byteLength(body)}},
        (res) => {
          const chunks: Buffer[] = [];
          res.on('data', d => chunks.push(d));
          res.on('end', () => {
            try { resolve(JSON.parse(Buffer.concat(chunks).toString('utf8'))); }
            catch (e) { reject(e); }
          });
        }
      );
      req.on('error', reject);
      req.write(body); req.end();
    });
  }
  dispose(): void {}
}

export function resolveCliPath(cfg: vscode.WorkspaceConfiguration): string {
  let p = cfg.get<string>('cliPath') || 'eve';
  if (process.platform === 'win32' && !p.endsWith('.exe')) p += '.exe';
  return p;
}
