const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');

module.exports = {
    cache: {
        type: 'filesystem',
        buildDependencies: {
            config: [__filename]
        }
    },
    // 開発用の設定
    mode: 'development',
    // エントリポイントとなるコード
    entry: './src/index.tsx',
    // バンドル後の js ファイルの出力先
    output: {
        path: path.resolve(__dirname, 'dist'),
        filename: 'bundle.js'
    },
    // import 時に読み込むファイルの拡張子
    resolve: {
        extensions: ['.js', '.json', '.ts', '.tsx']
    },
    // ソースマップファイルの出力設定
    devtool: 'source-map',
    module: {
        rules: [
            // TypeScript ファイル (.ts/.tsx) を変換できるようにする
            {
                test: /\.ts(x?)$/,
                use: "ts-loader",
                include: path.resolve(__dirname, 'src'),
                exclude: /node_modules/
            }, {
                test: /\.js(x?)$/,
                loader: 'babel-loader',
                exclude: /node_modules/,
            },
        ]
    },
    plugins: [
        // HTML ファイルの出力設定
        new HtmlWebpackPlugin({
            template: './src/index.html'
        })
    ],
    devServer: {
        port: '3000',
        contentBase: "dist",
        hot: true,
        open: true
    },
};
