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
  return new Promise((resolve, reject) => {
    https.get(url, (response) => {
      if (response.statusCode === 302 || response.statusCode === 301) {
        download(response.headers.location).then(resolve).catch(reject);
        return;
      }

      if (response.statusCode !== 200) {
        reject(new Error(`Failed to download: HTTP ${response.statusCode} from ${url}`));
        return;
      }

      response.pipe(file);
      file.on('finish', () => {
        file.close();
        fs.chmodSync(binPath, 0o755);

        // Create symlink in npm bin directory since npm processes bin entries before postinstall runs
        const npmBin = process.env.npm_config_prefix
          ? path.join(process.env.npm_config_prefix, 'bin')
          : path.dirname(process.execPath);
        const symlinkPath = path.join(npmBin, `opork${ext}`);

        try {
          if (fs.existsSync(symlinkPath)) {
            fs.unlinkSync(symlinkPath);
          }
          fs.symlinkSync(binPath, symlinkPath);
          console.log('opork installed successfully!');
        } catch (symlinkErr) {
          // Symlink creation may fail due to permissions, but binary is still usable
          console.log('opork binary installed. You may need to add it to your PATH manually.');
          console.log(`Binary location: ${binPath}`);
        }
        resolve();
      });
    }).on('error', (err) => {
      fs.unlink(binPath, () => {});
      reject(err);
    });
  });
}

download(downloadUrl).catch((err) => {
  console.error(`Download failed: ${err.message}`);
  process.exit(1);
});
