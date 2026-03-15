import type { CapacitorConfig } from '@capacitor/cli';

const config: CapacitorConfig = {
    appId: 'com.example.app',
    appName: 'trip-tuner',
    webDir: 'dist/trip-tuner/browser',
    server: {
        cleartext: true
    },
};

export default config;
