import React from 'react';
import { FaSearch, FaSync, FaTv } from 'react-icons/fa';

function Header({ searchQuery, onSearch, onReload }) {
  return (
    <header className="app-header">
      <div className="header-brand">
        <div className="brand-icon">
          <FaTv />
        </div>
        <span className="brand-title">Online TV</span>
      </div>

      <div className="header-search">
        <div className="search-box">
          <FaSearch className="search-icon" />
          <input
            type="text"
            placeholder="Search channels..."
            value={searchQuery}
            onChange={(e) => onSearch(e.target.value)}
            className="search-input"
          />
        </div>
      </div>

      <div className="header-actions">
        <button className="btn-reload" onClick={onReload} title="Refresh playlist">
          <FaSync size={12} />
          Refresh
        </button>
      </div>
    </header>
  );
}

export default Header;
