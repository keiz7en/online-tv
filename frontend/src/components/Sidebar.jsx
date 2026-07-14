import React from 'react';
import { FaList, FaNewspaper, FaFutbol, FaFilm, FaMusic, FaGlobe, FaStar, FaChurch, FaBaby, FaDesktop } from 'react-icons/fa';

const categoryIcons = {
  'All': FaList,
  'Sports': FaFutbol,
  'News': FaNewspaper,
  'Movies': FaFilm,
  'Music': FaMusic,
  'Live': FaDesktop,
  'Bangla': FaStar,
  'Indian': FaStar,
  'Religious': FaChurch,
  'Kids': FaBaby,
  'Documentary': FaGlobe,
};

function Sidebar({ categories, selectedCategory, onSelectCategory, channelCount }) {
  return (
    <aside className="sidebar">
      <div className="sidebar-title">Categories</div>
      <nav className="category-nav">
        <button
          className={`category-item ${selectedCategory === 'All' ? 'active' : ''}`}
          onClick={() => onSelectCategory('All')}
        >
          <FaList className="category-icon" />
          <span className="category-name">All Channels</span>
          <span className="category-count">{channelCount}</span>
        </button>

        {categories && categories.map((cat) => {
          const Icon = categoryIcons[cat] || FaGlobe;
          return (
            <button
              key={cat}
              className={`category-item ${selectedCategory === cat ? 'active' : ''}`}
              onClick={() => onSelectCategory(cat)}
            >
              <Icon className="category-icon" />
              <span className="category-name">{cat}</span>
            </button>
          );
        })}
      </nav>
    </aside>
  );
}

export default Sidebar;
