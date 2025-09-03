import { Component, ChangeDetectionStrategy, input, signal, OnInit, effect, viewChild, ElementRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { File } from '../../models/message.model';

interface TreeNode {
  name: string;
  type: 'file' | 'folder';
  data?: File;
  children?: TreeNode[];
  path: string;
}

declare var Prism: any;

@Component({
  selector: 'app-file-preview',
  templateUrl: './file-preview.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [CommonModule],
})
export class FilePreviewComponent implements OnInit {
  files = input.required<File[]>();

  fileTree = signal<TreeNode[]>([]);
  openTabs = signal<File[]>([]);
  activeFile = signal<File | null>(null);
  expandedFolders = signal<Set<string>>(new Set());

  private readonly codeBlock = viewChild<ElementRef<HTMLElement>>('codeBlock');

  constructor() {
    effect(() => {
      const codeElement = this.codeBlock()?.nativeElement;
      if (this.activeFile() && codeElement) {
        // Use a timeout to ensure the view is updated with new content before highlighting
        setTimeout(() => {
          Prism.highlightElement(codeElement);
        }, 0);
      }
    });
  }

  ngOnInit(): void {
    const tree = this.buildFileTree(this.files());
    this.fileTree.set(tree);

    // Automatically open the first file if available
    if (this.files().length > 0) {
      this.openFile(this.files()[0]);
    }
  }

  toggleFolder(path: string): void {
    this.expandedFolders.update(currentSet => {
      const newSet = new Set(currentSet);
      if (newSet.has(path)) {
        newSet.delete(path);
      } else {
        newSet.add(path);
      }
      return newSet;
    });
  }

  openFile(file: File): void {
    if (!this.openTabs().some(tab => tab.path === file.path)) {
      this.openTabs.update(tabs => [...tabs, file]);
    }
    this.setActiveFile(file);
  }
  
  selectTab(tab: File): void {
    this.setActiveFile(tab);
  }

  closeTab(tabToClose: File, event: MouseEvent): void {
    event.stopPropagation(); // Prevent selecting tab when clicking close

    const tabs = this.openTabs();
    const index = tabs.findIndex(tab => tab.path === tabToClose.path);

    if (index > -1) {
      const newTabs = tabs.filter(tab => tab.path !== tabToClose.path);
      this.openTabs.set(newTabs);

      if (this.activeFile()?.path === tabToClose.path) {
        if (newTabs.length > 0) {
          const newActiveIndex = Math.max(0, index - 1);
          this.setActiveFile(newTabs[newActiveIndex]);
        } else {
          this.setActiveFile(null);
        }
      }
    }
  }

  private setActiveFile(file: File | null): void {
    this.activeFile.set(file);

    if (file) {
      // Auto-expand folders to reveal the selected file
      this.expandedFolders.update(currentSet => {
        const newSet = new Set(currentSet);
        const pathParts = file.path.split('/');
        
        let currentPath = '';
        // Iterate over path parts to build and add parent folder paths
        for (let i = 0; i < pathParts.length - 1; i++) {
          currentPath = currentPath ? `${currentPath}/${pathParts[i]}` : pathParts[i];
          newSet.add(currentPath);
        }
        return newSet;
      });
    }
  }

  private buildFileTree(files: File[]): TreeNode[] {
    const fileTree: { [key: string]: any } = {};
    const expanded = new Set<string>();

    files.forEach(file => {
      let level = fileTree;
      const pathParts = file.path.split('/');
      pathParts.forEach((part, index) => {
        const currentPath = pathParts.slice(0, index + 1).join('/');

        if (index === pathParts.length - 1) {
          // This part is a file
          level[part] = { name: part, type: 'file', data: file, path: file.path };
        } else {
          // This part is a folder
          if (!level[part]) {
            level[part] = { name: part, type: 'folder', children: {}, path: currentPath };
            // Auto-expand root folders by default
            if (index === 0) {
                expanded.add(currentPath);
            }
          }
          level = level[part].children;
        }
      });
    });

    this.expandedFolders.set(expanded);

    const convertTreeToArray = (tree: { [key: string]: any }): TreeNode[] => {
      return Object.values(tree).map((node: any) => {
        if (node.type === 'folder') {
          node.children = convertTreeToArray(node.children);
        }
        return node;
      }).sort((a, b) => {
        if (a.type === 'folder' && b.type === 'file') return -1;
        if (a.type === 'file' && b.type === 'folder') return 1;
        return a.name.localeCompare(b.name);
      });
    };
    
    return convertTreeToArray(fileTree);
  }
}