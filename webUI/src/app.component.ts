import { Component, ChangeDetectionStrategy, signal, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ChatComponent } from './components/chat/chat.component';
import { SidePanelComponent } from './components/side-panel/side-panel.component';
import { SplashScreenComponent } from './components/splash-screen/splash-screen.component';

@Component({
  selector: 'app-root',
  template: `
    @if (showSplash()) {
      <div class="w-full h-screen animated-gradient">
        <app-splash-screen></app-splash-screen>
      </div>
    } @else {
      <main class="w-full h-screen animated-gradient flex flex-col relative overflow-hidden animate-fade-in">
        <app-chat class="flex-grow min-h-0 flex flex-col"></app-chat>
        <app-side-panel></app-side-panel>
      </main>
    }
  `,
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [CommonModule, ChatComponent, SidePanelComponent, SplashScreenComponent],
})
export class AppComponent implements OnInit {
  showSplash = signal(true);

  ngOnInit(): void {
    setTimeout(() => {
      this.showSplash.set(false);
    }, 4000); // Time for the splash screen text animation to complete
  }
}
