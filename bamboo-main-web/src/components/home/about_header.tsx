/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋(https://www.x-lf.com)
 *
 * 本文件包含 XiaoMain 的源代码，该项目的所有源代码均遵循MIT开源许可证协议。
 * --------------------------------------------------------------------------------
 * 许可证声明：
 *
 * 版权所有 (c) 2016-2024 筱锋。保留所有权利。
 *
 * 本软件是“按原样”提供的，没有任何形式的明示或暗示的保证，包括但不限于
 * 对适销性、特定用途的适用性和非侵权性的暗示保证。在任何情况下，
 * 作者或版权持有人均不承担因软件或软件的使用或其他交易而产生的、
 * 由此引起的或以任何方式与此软件有关的任何索赔、损害或其他责任。
 *
 * 使用本软件即表示您了解此声明并同意其条款。
 *
 * 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 * 免责声明：
 *
 * 使用本软件的风险由用户自担。作者或版权持有人在法律允许的最大范围内，
 * 对因使用本软件内容而导致的任何直接或间接的损失不承担任何责任。
 * --------------------------------------------------------------------------------
 */

import {Link, useLocation} from "react-router-dom";
import {AppstoreAddOutlined, HomeOutlined, LinkOutlined, UserOutlined} from "@ant-design/icons";

export function AboutHeaderComponent() {
    const location = useLocation();

    function SelectColorForPhone(getLocation: string): string {
        if (location.pathname.includes(getLocation)) {
            return "bg-sky-100 p-2 text-sm font-medium text-sky-600"
        } else {
            return "p-2 text-sm font-medium text-gray-500 hover:bg-gray-50 hover:text-gray-700"
        }
    }

    function SelectColorForDesk(getLocation: string): string {
        if (location.pathname === getLocation) {
            return "text-blue-600 hover:text-blue-600/75"
        } else {
            return "text-gray-500 hover:text-gray-500/75"
        }
    }

    return (
        <>
            <header className="hidden md:block">
                <div className="mx-auto max-w-screen-xl px-4 sm:px-6 lg:px-8">
                    <div className="flex md:h-14 lg:h-16 items-center justify-between">
                        <div className="flex-1 md:flex md:items-center md:gap-12">
                            <Link className="block text-teal-600" to={"/"}>
                                <span className="sr-only">Home</span>
                                <img src={"/favicon.png"} alt={"WebLogo"} className={"size-9"}/>
                            </Link>
                        </div>
                        <div className="flex items-center md:gap-12">
                            <nav aria-label="Global" className="hidden md:block">
                                <ul className="flex items-center gap-3 text-sm font-medium">
                                    <li>
                                        <Link className={`transition py-1 px-3 ${SelectColorForDesk("/")}`}
                                              to={"/"}>首页</Link>
                                    </li>
                                    <li>
                                        <Link className={`transition py-1 px-3 ${SelectColorForDesk("/me/about")}`}
                                              to="/me/about">关于我</Link>
                                    </li>
                                    <li>
                                        <Link className={`transition py-1 px-3 ${SelectColorForDesk("/me/friends")}`}
                                              to="/me/friends">友链</Link>
                                    </li>
                                    <li>
                                        <Link className={`transition py-1 px-3 ${SelectColorForDesk("/me/sponsor")}`}
                                              to="/me/sponsor">赞助</Link>
                                    </li>
                                </ul>
                            </nav>
                        </div>
                    </div>
                </div>
            </header>
            {/* 手机菜单样式 */}
            <div className={"fixed md:hidden w-lvw bottom-0 left-0 p-3 bg-white"}>
                <div className={"flex items-center justify-center"}>
                    <nav className="grid grid-cols-4 gap-3 w-full" aria-label="Tabs">
                        <Link to={"/"}
                              className={`transition shrink-0 rounded-lg flex items-center justify-center gap-1 text-gray-500 hover:text-gray-500/75`}>
                            <HomeOutlined/>
                            <span>首页</span>
                        </Link>
                        <Link to={"/me/about"}
                              className={`transition shrink-0 rounded-lg flex items-center justify-center gap-1 ${SelectColorForPhone("about")}`}>
                            <UserOutlined/>
                            <span>是我</span>
                        </Link>
                        <Link to={"/me/friends"}
                              className={`transition shrink-0 rounded-lg flex items-center justify-center gap-1 ${SelectColorForPhone("friends")}`}>
                            <LinkOutlined/>
                            <span>友链</span>
                        </Link>
                        <Link to={"/me/other"}
                              className={`transition shrink-0 rounded-lg flex items-center justify-center gap-1 ${SelectColorForPhone("other")}`}>
                            <AppstoreAddOutlined/>
                            <span>其他</span>
                        </Link>
                    </nav>
                </div>
            </div>
        </>
    );
}
