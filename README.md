# 🎙️ VoiceLine - Audio Transcription & Summarization with AWS Lambda

This project provides two AWS Lambda functions for end-to-end processing of voice input. It transcribes raw audio using OpenAI’s Whisper model and summarizes the resulting text into a structured format, ideal for analyzing sales conversations.

---

## 📦 Project Structure

### 1. `VoiceLineProcessAudioLambda`
**Purpose**: Transcribes user voice input.

- **Input**: Audio file (e.g., `.mp3`, `.wav`)
- **Process**: 
  - Sends audio to OpenAI's [Whisper](https://openai.com/research/whisper) model.
  - Receives and returns the text transcript.
- **Output**: Raw text transcript of spoken audio.
- [Endpoint](https://l2rulpi2rl.execute-api.eu-central-1.amazonaws.com/default/voicelinemain) (Post Request)

### 2. `VoiceLineSummarizeLambda`
**Purpose**: Summarizes the transcript into structured conversation data.

- **Input**: Transcript text (e.g., from the first Lambda).
- **Process**:
  - Sends the transcript to a GPT-based model (e.g., ChatGPT).
  - Extracts and formats key discussion points.
- **Output**: A structured summary (e.g., customer needs, objections, follow-ups).
- [Endpoint](https://l2rulpi2rl.execute-api.eu-central-1.amazonaws.com/default/voicelinetest?text="Hello") (Get Request)


---

## 🚀 Deployment

Both functions are deployed on AWS Lambda and also integrated with API Gateway.


