import React from 'react';
import { FaStop, FaPause, FaVolumeUp, FaVolumeMute, FaTv } from 'react-icons/fa';

function PlayerBar({ channel, playing, volume, onStop, onTogglePause, onVolumeChange }) {
  return (
    <div className="player-bar">
      <div className="player-info">
        <FaTv className="player-icon" />
        <div className="player-channel-info">
          <span className="player-channel-name">{channel?.name || 'No channel'}</span>
          <span className="player-channel-category">{channel?.category || ''}</span>
        </div>
      </div>

      <div className="player-controls">
        <button className="control-btn" onClick={onTogglePause} title="Pause">
          <FaPause size={12} />
        </button>
        <button className="control-btn stop-btn" onClick={onStop} title="Stop">
          <FaStop size={10} />
        </button>
      </div>

      <div className="player-volume">
        <button
          className="control-btn"
          onClick={() => onVolumeChange(volume > 0 ? 0 : 100)}
          title={volume > 0 ? 'Mute' : 'Unmute'}
          style={{ width: 32, height: 32 }}
        >
          {volume > 0 ? <FaVolumeUp size={12} /> : <FaVolumeMute size={12} />}
        </button>
        <input
          type="range"
          min="0"
          max="100"
          value={volume}
          onChange={(e) => onVolumeChange(parseInt(e.target.value))}
          className="volume-slider"
        />
        <span className="volume-value">{volume}%</span>
      </div>
    </div>
  );
}

export default PlayerBar;
