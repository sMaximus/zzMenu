/*
 * @Author: 张家铭 zhangjiaming@sursenelec.com
 * @Date: 2025-08-22 10:33:04
 * @LastEditors: 张家铭 zhangjiaming@sursenelec.com
 * @LastEditTime: 2025-08-28 11:23:21
 * @FilePath: \zzMenu\vite.config.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import { resolve } from "path";
import tailwindcss from "@tailwindcss/vite";
// import legacy from "@vitejs/plugin-legacy";
const rootDir = resolve(__dirname);
const srcDir = resolve(rootDir, "src");
const pageDir = resolve(srcDir, "pages");

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: {
      "@": srcDir,
      "@public": resolve(rootDir, "public"),
      "@pages": pageDir,
    },
  },
});
