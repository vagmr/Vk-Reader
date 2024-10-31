import React, { useState } from 'react';
import { SearchBooks } from '../../wailsjs/go/main/App';
import { BookCard } from './BookCard';
import {main} from "../../wailsjs/go/models"
import { LogPrint ,LogError} from '../../wailsjs/runtime/runtime';

export const BookSearch: React.FC = () => {
  const [keyword, setKeyword] = useState('');
  const [books, setBooks] = useState<main.Book[]>([]);
  const [loading, setLoading] = useState(false);

  const handleSearch = async () => {
    if (!keyword.trim()) return;
    
    setLoading(true);
    try {
      const results = await SearchBooks(keyword);
      setBooks(results);
      LogPrint("搜索成功");
    } catch (err) {
      console.error('搜索失败:', err);
      LogError(`搜索失败: ${err}`);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="book-search">
      <div className="search-bar">
        <input 
          type="text"
          value={keyword}
          onChange={(e) => setKeyword(e.target.value)}
          placeholder="输入小说名称或作者"
        />
        <button 
          onClick={handleSearch}
          disabled={loading}
        >
          {loading ? '搜索中...' : '搜索'}
        </button>
      </div>

      <div className="book-list">
        {books.map((book) => (
          <BookCard key={book.id} book={book} />
        ))}
      </div>
    </div>
  );
}; 