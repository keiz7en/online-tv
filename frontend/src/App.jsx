import React, { useState, useEffect } from 'react';
import { GetChannels, GetCategories, SearchChannels, PlayChannel, StopPlayback, TogglePause, SetVolume, GetVolume, IsPlaying, ReloadPlaylist } from '../wailsjs/go/main/App';
import Header from './components/Header';
import Sidebar from './components/Sidebar';
import ChannelGrid from './components/ChannelGrid';
import PlayerBar from './components/PlayerBar';
import './App.css';

function App() {
  const [channels, setChannels] = useState([]);
  const [categories, setCategories] = useState([]);
  const [selectedCategory, setSelectedCategory] = useState('All');
  const [searchQuery, setSearchQuery] = useState('');
  const [playing, setPlaying] = useState(false);
  const [currentChannel, setCurrentChannel] = useState(null);
  const [volume, setVolumeState] = useState(100);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadPlaylist();
  }, []);

  useEffect(() => {
    loadChannels();
  }, [selectedCategory, searchQuery]);

  const loadPlaylist = async () => {
    setLoading(true);
    try {
      await ReloadPlaylist();
      await loadChannels();
    } catch (err) {
      console.error('Failed to load playlist:', err);
    }
    setLoading(false);
  };

  const loadChannels = async () => {
    try {
      let data;
      if (searchQuery) {
        data = await SearchChannels(searchQuery);
      } else {
        data = await GetChannels();
      }
      setChannels(data || []);

      const cats = await GetCategories();
      setCategories(cats || []);
    } catch (err) {
      console.error('Failed to load channels:', err);
    }
  };

  const handleSearch = (query) => {
    setSearchQuery(query);
  };

  const handleCategorySelect = (category) => {
    setSelectedCategory(category);
    setSearchQuery('');
  };

  const handlePlayChannel = async (channel) => {
    try {
      await PlayChannel(channel.url);
      setPlaying(true);
      setCurrentChannel(channel);
    } catch (err) {
      console.error('Failed to play channel:', err);
    }
  };

  const handleStop = async () => {
    try {
      await StopPlayback();
      setPlaying(false);
      setCurrentChannel(null);
    } catch (err) {
      console.error('Failed to stop:', err);
    }
  };

  const handleTogglePause = async () => {
    try {
      await TogglePause();
    } catch (err) {
      console.error('Failed to toggle pause:', err);
    }
  };

  const handleVolumeChange = async (newVol) => {
    try {
      await SetVolume(newVol);
      setVolumeState(newVol);
    } catch (err) {
      console.error('Failed to set volume:', err);
    }
  };

  const filteredChannels = selectedCategory === 'All'
    ? channels
    : channels.filter(ch => ch.category === selectedCategory);

  return (
    <div id="app" className="app-container">
      <Header
        searchQuery={searchQuery}
        onSearch={handleSearch}
        onReload={loadPlaylist}
      />

      <div className="main-content">
        <Sidebar
          categories={categories}
          selectedCategory={selectedCategory}
          onSelectCategory={handleCategorySelect}
          channelCount={channels.length}
        />

        <ChannelGrid
          channels={filteredChannels}
          onPlayChannel={handlePlayChannel}
          currentChannel={currentChannel}
          loading={loading}
        />
      </div>

      {playing && (
        <PlayerBar
          channel={currentChannel}
          playing={playing}
          volume={volume}
          onStop={handleStop}
          onTogglePause={handleTogglePause}
          onVolumeChange={handleVolumeChange}
        />
      )}
    </div>
  );
}

export default App;
