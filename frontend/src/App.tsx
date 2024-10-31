import { BookSearch } from './components/BookSearch';
import { Settings } from './components/Settings';
import { DownloadProgress } from './components/DownloadProgress';
import { GetCookie } from '../wailsjs/go/main/App'
import './App.css';
import { useEffect } from 'react';

function App() {
  useEffect(() => {
    GetCookie()
  }, [])
  return (
    <div className="container">
      <header className="header">
        <h1>番茄小说下载器</h1>
      </header>
      
      <main className="reader-container">
        <div className="content">
          {/* 搜索组件 */}
          <BookSearch />
          
          {/* 下载进度组件 */}
          <DownloadProgress />
          
          {/* 设置组件 */}
          <Settings />
        </div>
      </main>
    </div>
  );
}

export default App;
