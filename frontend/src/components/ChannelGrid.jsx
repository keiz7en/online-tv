import React from 'react';
import ChannelCard from './ChannelCard';

function ChannelGrid({ channels, onPlayChannel, currentChannel, loading }) {
  if (loading) {
    return (
      <div className="channel-grid-container">
        <div className="loading-state">
          <div className="loading-spinner"></div>
          <span>Loading channels...</span>
        </div>
      </div>
    );
  }

  if (!channels || channels.length === 0) {
    return (
      <div className="channel-grid-container">
        <div className="empty-state">
          <span>No channels found</span>
        </div>
      </div>
    );
  }

  return (
    <div className="channel-grid-container">
      <div className="channel-grid">
        {channels.map((channel, index) => (
          <ChannelCard
            key={`${channel.url}-${index}`}
            channel={channel}
            isPlaying={currentChannel?.url === channel.url}
            onPlay={() => onPlayChannel(channel)}
          />
        ))}
      </div>
    </div>
  );
}

export default ChannelGrid;
