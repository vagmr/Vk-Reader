import React, { useState } from 'react';
import { DownloadBook } from '../../wailsjs/go/main/App';
import { main } from '../../wailsjs/go/models';

interface BookProps {
    book: main.Book;
}

export const BookCard: React.FC<BookProps> = ({ book }) => {
  const [downloading, setDownloading] = useState(false);

  const handleDownload = async () => {
    setDownloading(true);
    try {
      await DownloadBook(book.id, 'txt'); // 默认txt格式
    } catch (err) {
      console.error('下载失败:', err);
    } finally {
      setDownloading(false);
    }
  };

  return (
    <div className="book-card">
      <img src={book.cover} alt={book.title} className="book-cover" />
      <div className="book-info">
        <h3>{book.title}</h3>
        <p>作者: {book.author}</p>
        <p>字数: {book.wordCount}</p>
        <p>状态: {book.status}</p>
        <button 
          onClick={handleDownload}
          disabled={downloading}
        >
          {downloading ? '下载中...' : '下载'}
        </button>
      </div>
    </div>
  );
}; 