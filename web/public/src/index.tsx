import React from 'react';
import ReactDOM from 'react-dom/client'; // 新しいインポートパス
import App from './component/App';

// コンテナを取得
const container = document.getElementById('app');

if (container !== null) {
    // createRootを使用してルートを作成（nullでないことを確認）
    const root = ReactDOM.createRoot(container);
    // ルートにAppコンポーネントをレンダリング
    root.render(<App id={1} ids={[1, 2, 3]} />);
} else {
    console.error('Failed to find the root element');
}