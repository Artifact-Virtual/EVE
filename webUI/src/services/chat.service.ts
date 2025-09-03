import { Injectable } from '@angular/core';
import { Message } from '../models/message.model';

@Injectable({
  providedIn: 'root',
})
export class ChatService {
  /**
   * This method sends the user's prompt to the Go backend and returns the assistant's response.
   */
  async getAssistantResponse(prompt: string): Promise<Message> {
    console.log('Sending prompt to Go backend:', prompt);
    
    try {
      const response = await fetch('/api/chat', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ prompt }),
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const assistantMessage: Message = await response.json();
      return assistantMessage;

    } catch (error) {
      console.error('Error fetching from Go backend:', error);
      // Return a user-friendly error message to be displayed in the chat
      return {
        sender: 'assistant',
        parts: [{
          type: 'text',
          content: 'Sorry, I couldn\'t connect to the backend. Please ensure the server is running.'
        }]
      };
    }
  }
}
