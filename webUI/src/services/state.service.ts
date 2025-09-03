import { Injectable, signal } from '@angular/core';
import { File } from '../models/message.model';

@Injectable({
  providedIn: 'root',
})
export class StateService {
  public readonly files = signal<File[]>([]);
  public readonly isPanelOpen = signal<boolean>(false);

  togglePanel(): void {
    this.isPanelOpen.update(isOpen => !isOpen);
  }
}
