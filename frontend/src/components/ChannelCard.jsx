import React from 'react';
import { FaPlay, FaTv } from 'react-icons/fa';

function ChannelCard({ channel, isPlaying, onPlay }) {
  return (
    <div
      className={`channel-card ${isPlaying ? 'playing' : ''}`}
      onClick={onPlay}
    >
      <div className="channel-logo">
        {channel.logo ? (
          <img
            src={channel.logo}
            alt={channel.name}
            loading="lazy"
            onError={(e) => {
              e.target.style.display = 'none';
              e.target.nextSibling.style.display = 'flex';
            }}
          />
        ) : null}
        <div className="channel-logo-fallback" style={{ display: channel.logo ? 'none' : 'flex' }}>
          <FaTv />
        </div>
        {isPlaying && (
          <div className="playing-indicator">
            <FaPlay size={16} />
          </div>
        )}
      </div>
      <div className="channel-info">
        <div className="channel-name">{channel.name}</div>
        <div className="channel-category">{channel.category}</div>
      </div>
    </div>
  );
}

export default ChannelCard;
