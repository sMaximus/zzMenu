/*
 * @Author: 张家铭 zhangjiaming@sursenelec.com
 * @Date: 2025-08-22 10:33:04
 * @LastEditors: 张家铭 zhangjiaming@sursenelec.com
 * @LastEditTime: 2025-08-22 18:13:00
 * @FilePath: \zzMenu\src\App.tsx
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { RouterProvider } from "react-router-dom";
import router from "./routers";
import { useState } from "react";
import "./App.css";

function App() {
  const [count, setCount] = useState(0);

  return (
    <>
      <RouterProvider router={router} />
    </>
  );
}

export default App;
