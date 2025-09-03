export interface File {
  path: string;
  content: string;
  language: string;
}

export interface TextPart {
  type: 'text';
  content: string;
}

export interface CodePart {
  type: 'code';
  content: string;
  language?: string;
}

export interface FilesPart {
  type: 'files';
  content: File[];
}

export type MessagePart = TextPart | CodePart | FilesPart;

export interface Message {
  sender: 'user' | 'assistant';
  parts: MessagePart[];
}
