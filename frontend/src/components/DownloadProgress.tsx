import React, { useEffect, useState } from 'react';
import { EventsOn } from '../../wailsjs/runtime/runtime';

interface ProgressData {
    total: number;
    current: number;
    chapterTitle: string;
    status: 'downloading' | 'completed' | 'error';
}

export const DownloadProgress: React.FC = () => {
    const [progress, setProgress] = useState<ProgressData | null>(null);

    useEffect(() => {
        const unsubscribe = EventsOn('download:progress', (data: ProgressData) => {
            setProgress(data);
        });

        return () => {
            unsubscribe();
        };
    }, []);

    if (!progress) return null;

    const percentage = Math.round((progress.current / progress.total) * 100);

    return (
        <div className="download-progress">
            <div className="progress-bar">
                <div 
                    className="progress-fill"
                    style={{ width: `${percentage}%` }}
                />
            </div>
            <div className="progress-info">
                <span>{progress.chapterTitle}</span>
                <span>{percentage}% ({progress.current}/{progress.total})</span>
            </div>
            {progress.status === 'error' && (
                <div className="error-message">
                    下载出错，请重试
                </div>
            )}
        </div>
    );
}; 