/*
 * @Author: 张家铭 zhangjiaming@sursenelec.com
 * @Date: 2025-08-22 17:25:08
 * @LastEditors: 张家铭 zhangjiaming@sursenelec.com
 * @LastEditTime: 2025-08-28 11:18:43
 * @FilePath: \zzMenu\src\routers\index.tsx
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { createBrowserRouter } from "react-router";
import Layout from "@/layout";

const router = [
  {
    path: "/",
    element: <Layout />,
    children: [],
  },
  // {
  //   path: "/login",
  //   element: <Login />,
  // },
  // {
  //   path: "/error",
  //   element: <Error />,
  // },
  // {
  //   path: "*",
  //   element: <NotFound />,
  // },
];

const routers = createBrowserRouter(router);

export { router };
export default routers;
