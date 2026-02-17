# Flutter 本地 → 自动连接 localhost:8282
cd frontend/flutter && flutter run -d chrome

# Vben Admin 本地 → vite proxy 代理到 localhost:8282
cd frontend/vben-admin && pnpm dev:naive

# 真机
open /Users/mmx/Documents/work/Github/e-woms/frontend/flutter/ios/Runner.xcworkspace

# 打包真机
flutter run --release --dart-define=ENV=prod -d 00008120-001220362E13C01E


# 打包到浏览器
flutter run -d chrome --web-port=8090 --web-browser-flag="--user-data-dir=$(pwd)/.chrome-profile"

# 构建好，再通过 Xcode Archive 分发
cd /Users/mmx/Documents/work/Github/e-woms/frontend/flutter
flutter build ios --release --dart-define=ENV=prod
# open Runner.xcworkspace → Product → Archive
