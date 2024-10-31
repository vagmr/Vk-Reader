import React, { useState, useEffect } from 'react';
import { GetConfig, SaveConfig } from '../../wailsjs/go/main/App';
import { main } from '../../wailsjs/go/models';

export const Settings: React.FC = () => {
  const [config, setConfig] = useState<main.Config>(main.Config.createFrom({
    savePath: '',
    saveFormat: 'txt',
    threadCount: 1,
    convertType: {
      txt: 'txt',
      epub: 'epub',
      html: 'html'
    },
    delay: {
      min: 50,
      max: 150
    }
  }));

  useEffect(() => {
    loadConfig();
  }, []);

  const loadConfig = async () => {
    try {
      const cfg = await GetConfig();
      setConfig(cfg);
    } catch (err) {
      console.error('加载配置失败:', err);
    }
  };

  const handleSave = async () => {
    try {
      await SaveConfig(config);
      console.log('配置已保存');
    } catch (err) {
      console.error('保存配置失败:', err);
    }
  };

  const handleChange = (field: keyof main.Config, value: any) => {
    setConfig(prev => main.Config.createFrom({
      ...prev,
      [field]: value
    }));
  };

  const handleDelayChange = (field: 'min' | 'max', value: number) => {
    setConfig(prev => main.Config.createFrom({
      ...prev,
      delay: {
        ...prev.delay,
        [field]: value
      }
    }));
  };

  return (
    <div className="settings">
      <h2>设置</h2>
      <div className="setting-item">
        <label>保存路径:</label>
        <input 
          type="text" 
          value={config.savePath}
          onChange={(e) => handleChange('savePath', e.target.value)}
        />
      </div>
      <div className="setting-item">
        <label>保存格式:</label>
        <select 
          value={config.saveFormat}
          onChange={(e) => handleChange('saveFormat', e.target.value)}
        >
          {Object.keys(config.convertType).map(format => (
            <option key={format} value={format}>
              {format.toUpperCase()}
            </option>
          ))}
        </select>
      </div>
      <div className="setting-item">
        <label>下载线程数:</label>
        <input 
          type="number"
          value={config.threadCount}
          onChange={(e) => handleChange('threadCount', parseInt(e.target.value))}
          min="1"
          max="10"
        />
      </div>
      <div className="setting-item">
        <label>下载延迟(毫秒):</label>
        <div className="delay-inputs">
          <input
            type="number"
            value={config.delay.min}
            onChange={(e) => handleDelayChange('min', parseInt(e.target.value))}
            min="0"
          />
          <span>-</span>
          <input
            type="number"
            value={config.delay.max}
            onChange={(e) => handleDelayChange('max', parseInt(e.target.value))}
            min={config.delay.min}
          />
        </div>
      </div>
      <button onClick={handleSave}>保存设置</button>
    </div>
  );
}; 