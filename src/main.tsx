/*
 * @Author: 张家铭 zhangjiaming@sursenelec.com
 * @Date: 2025-08-22 10:33:04
 * @LastEditors: 张家铭 zhangjiaming@sursenelec.com
 * @LastEditTime: 2025-08-22 18:09:15
 * @FilePath: \zzMenu\src\main.tsx
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./index.css";
import App from "./App.tsx";

createRoot(document.getElementById("root")!).render(<App />);
