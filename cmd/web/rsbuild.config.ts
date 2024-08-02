import {defineConfig} from '@rsbuild/core';
import {pluginReact} from '@rsbuild/plugin-react';
import {pluginSass} from "@rsbuild/plugin-sass";
// @ts-expect-error complains that can not load json but loads it fine
import {dependencies} from './package.json';
import CompressionPlugin from 'compression-webpack-plugin/dist';

// Customize these variables to match your app requirements
const regex = /-/g;
const name: string = 'example-app'.replace(regex, '_'); // Should be globally unique
const serverPort: number = Number('3001'); // Should be globally unique with a group eg. (next, cloud)

export default defineConfig({
    plugins: [pluginReact(), pluginSass()],
    html: {
      crossorigin: 'anonymous',
    },
    server: {
        open: false,
        port: serverPort,
        host: '0.0.0.0',
        compress: true,
        headers: {
            'Access-Control-Allow-Origin': '*',
        },
    },
    dev: {
        assetPrefix: 'http://localhost:3001/',
    },
    output: {
      minify: true
    },
    moduleFederation: {
        options: {
            name: name,
            filename: 'remoteEntry.js',
            exposes: {
                "./App": './src/components/Counter.tsx'
            },
            shared: {
                ...dependencies,
                'react': {
                    requiredVersion: dependencies['react'],
                    singleton: true
                },
                'react-dom': {
                    requiredVersion: dependencies['react-dom'],
                    singleton: true
                },
                'react-router-dom': {
                    requiredVersion: dependencies['react-router-dom'],
                    singleton: true
                },
            },
        }
    },
    tools: {
        rspack: (config, {appendPlugins}) => {
            // You need to set a unique value that is not equal to other applications
            config.output!.uniqueName = name;
            config.output!.publicPath = "auto";

            appendPlugins([
                new CompressionPlugin(),
            ]);
        },
    },
});
