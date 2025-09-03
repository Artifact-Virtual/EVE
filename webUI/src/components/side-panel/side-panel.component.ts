import { Component, ChangeDetectionStrategy, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { StateService } from '../../services/state.service';
import { FilePreviewComponent } from '../file-preview/file-preview.component';
import { File } from '../../models/message.model';
import { Signal } from '@angular/core';

@Component({
  selector: 'app-side-panel',
  templateUrl: './side-panel.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [CommonModule, FilePreviewComponent],
})
export class SidePanelComponent {
  private stateService = inject(StateService);

  isPanelOpen: Signal<boolean> = this.stateService.isPanelOpen;
  files: Signal<File[]> = this.stateService.files;
  
  activeTab = signal('files');

  closePanel(): void {
    this.stateService.isPanelOpen.set(false);
  }
}
