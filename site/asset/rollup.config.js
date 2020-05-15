import nodeResolve from 'rollup-plugin-node-resolve';
import commonjs from 'rollup-plugin-commonjs';
import { terser } from 'rollup-plugin-terser';
import json from '@rollup/plugin-json';

export default {
    input: {
        'main': 'dev/javascript/main.js',
    },
    output: {
        dir: 'live/javascript',
        format: 'es',
        globals: [],
    },
    plugins: [
        json(),
        nodeResolve(),
        commonjs(),
        terser()
    ],
    external: []
}