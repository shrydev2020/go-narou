const { build } = require('esbuild')
const path = require('path')
// optionsの定義
const options = {
    // 以下のdefineプロパティを設定しない場合Reactのプロジェクトの実行時にエラーが出ます
    define: { 'process.env.NODE_ENV': process.env.NODE_ENV },
    entryPoints: [path.resolve(__dirname, 'src/index.tsx')],
    bundle: true,
    target: 'es2016',
    platform: 'browser',
    outdir: path.resolve(__dirname, 'dist'),
    tsconfig: path.resolve(__dirname, 'tsconfig.json')
}
// Buildの実行
// @ts-ignore
build(options).catch((err) => {
    process.stderr.write(err.stderr)
    process.exit(1)
})