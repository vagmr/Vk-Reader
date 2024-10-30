import React from 'react';
import './ReaderView.css';

interface ReaderViewProps {
  // 属性定义
}

const ReaderView: React.FC<ReaderViewProps> = () => {
  return (
    <div className="reader-container">
      <div className="reader-content">
        {/* 阅读内容区域 */}
      </div>
      <div className="reader-controls">
        {/* 控制按钮区域 */}
      </div>
    </div>
  );
};

export default ReaderView; 