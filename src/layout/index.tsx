/*
 * @Author: 张家铭 zhangjiaming@sursenelec.com
 * @Date: 2025-08-22 17:30:53
 * @LastEditors: 张家铭 zhangjiaming@sursenelec.com
 * @LastEditTime: 2025-08-28 16:18:36
 * @FilePath: \zzMenu\src\layout\index.tsx
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

import React from "react";
import { SideBar, Tabs } from "antd-mobile";
import header from "@/assets/images/header.png";

const tabs = [
  {
    key: "1",
    title: "大厨",
  },
  {
    key: "2",
    title: "家常菜",
  },
  {
    key: "3",
    title: "主食",
  },
  {
    key: "4",
    title: "小吃",
  },
  // {
  //   key: "5",
  //   title: "音频",
  // },
];

const menuList = [
  {
    key: "1",
    title: "菜单",
  },
  {
    key: "2",
    title: "订单",
  },
  {
    key: "3",
    title: "我的",
  },
];

const layoutComponent = () => {
  return (
    <div className="h-screen flex flex-col">
      <img src={header} alt="header" className="w-full" />

      <div className="flex flex-1 flex-col">
        <SideBar>
          {tabs.map((item) => (
            <SideBar.Item key={item.key} title={item.title} />
          ))}
        </SideBar>

        <Tabs>
          {menuList.map((item) => (
            <Tabs.Tab key={item.key} title={item.title} />
          ))}
        </Tabs>
      </div>
    </div>
  );
};

export default layoutComponent;
