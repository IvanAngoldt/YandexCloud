import React, { useState, useRef } from 'react';
import { sendAudioToServer } from '../../api/api';
import '../../styles/Blog/MicrophoneButton.css'; // Подключение стилей кнопки записи

const MicrophoneButton = ({ onResult, onWaiting }) => {
  const [recording, setRecording] = useState(false);
  const mediaRecorderRef = useRef(null);
  const audioChunks = useRef([]);

  const startRecording = async () => {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
      const mediaRecorder = new MediaRecorder(stream, { mimeType: 'audio/ogg; codecs=opus' });

      mediaRecorderRef.current = mediaRecorder;
      audioChunks.current = [];
      onWaiting(true);

      mediaRecorder.ondataavailable = (event) => {
        audioChunks.current.push(event.data);
      };

      mediaRecorder.onstop = async () => {
        const audioBlob = new Blob(audioChunks.current, { type: 'audio/ogg; codecs=opus' });
        await handleSendAudio(audioBlob);
      };

      mediaRecorder.start();
      setRecording(true);
    } catch (error) {
      console.error('Ошибка при доступе к микрофону:', error);
    }
  };

  const stopRecording = () => {
    if (mediaRecorderRef.current) {
      mediaRecorderRef.current.stop();
      setRecording(false);
    }
  };

  const handleSendAudio = async (audioBlob) => {
    try {
      const result = await sendAudioToServer(audioBlob);
      onResult(result.result || '');
    } catch (error) {
      onResult('Произошла ошибка при распознавании.');
    } finally {
      onWaiting(false);
    }
  };

  return (
    <button
      type="button" // Добавлен тип кнопки
      className={`mic-button ${recording ? 'recording' : ''}`}
      onClick={recording ? stopRecording : startRecording}
      title={recording ? 'Остановить запись' : 'Начать запись'}
    >
      {recording ? '●' : '🎤'}
    </button>
  );
};

export default MicrophoneButton;
