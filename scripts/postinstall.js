#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const https = require('https');
const { execSync } = require('child_process');

const pkg = require('../package.json');
const version = pkg.version;

const platformMap = {
  darwin: 'darwin',
  linux: 'linux',
  win32: 'windows',
};

const archMap = {
  x64: 'amd64',
  arm64: 'arm64',
};

const platform = platformMap[process.platform];
const arch = archMap[process.arch];

if (!platform || !arch) {
  console.error(`Unsupported platform: ${process.platform}-${process.arch}`);
  process.exit(1);
}

const ext = platform === 'windows' ? '.exe' : '';
const binaryName = `opork-${platform}-${arch}${ext}`;
const downloadUrl = `https://github.com/OverseedAI/overpork/releases/download/v${version}/${binaryName}`;

const binDir = path.join(__dirname, '..', 'bin');
const binPath = path.join(binDir, `opork${ext}`);

if (!fs.existsSync(binDir)) {
  fs.mkdirSync(binDir, { recursive: true });
}

console.log(`Downloading opork v${version} for ${platform}-${arch}...`);

const file = fs.createWriteStream(binPath);

function download(url) {
  https.get(url, (response) => {
    if (response.statusCode === 302 || response.statusCode === 301) {
      download(response.headers.location);
      return;
    }

    if (response.statusCode !== 200) {
      console.error(`Failed to download: HTTP ${response.statusCode}`);
      console.error(`URL: ${url}`);
      process.exit(1);
    }

    response.pipe(file);
    file.on('finish', () => {
      file.close();
      fs.chmodSync(binPath, 0o755);
      console.log('opork installed successfully!');
    });
  }).on('error', (err) => {
    fs.unlink(binPath, () => {});
    console.error(`Download failed: ${err.message}`);
    process.exit(1);
  });
}

download(downloadUrl);
