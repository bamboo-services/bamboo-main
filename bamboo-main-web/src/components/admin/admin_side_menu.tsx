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

import {
    BarsOutlined,
    CaretDownOutlined,
    HomeOutlined,
    LinkOutlined,
    SettingOutlined,
    UsergroupAddOutlined
} from "@ant-design/icons";
import {UserCurrentEntity} from "../../resources/ts/models/entity/user_current_entity.ts";
import {Link, useLocation} from "react-router-dom";

export function AdminSideMenuComponent(userCurrent: UserCurrentEntity) {
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
            return "bg-gray-100 font-medium text-gray-700"
        } else {
            return "text-gray-500 hover:bg-gray-100 hover:text-gray-700"
        }
    }

    return (
        <>
            {/* 电脑菜单样式 */}
            <div className="hidden md:flex h-screen flex-col justify-between border-e bg-white fixed left-0 top-0 w-[200px] lg:w-[250px]">
                <div className="px-4 py-6">
                    <Link to={"/"} className="grid h-10 place-content-center rounded-lg bg-gray-100 text-xs text-gray-600">
                      Logo
                    </Link>
                    <ul className="mt-6 space-y-1">
                        <li>
                            <Link to={"/admin/dashboard"}
                                  className={`transition flex items-center gap-2 rounded-lg px-4 py-2 text-sm ${SelectColorForDesk("/admin/dashboard")}`}>
                                <HomeOutlined/>
                                <span>首页</span>
                            </Link>
                        </li>
                        <li>
                            <details className="group">
                                <summary
                                    className="transition flex cursor-pointer items-center justify-between rounded-lg px-4 py-2 text-gray-500 hover:bg-gray-100 hover:text-gray-700"
                                >
                                    <span className="text-sm font-medium flex items-center gap-2">
                                        <UsergroupAddOutlined />
                                        <span>友链</span>
                                    </span>
                                    <span className="shrink-0 transition duration-300 group-open:-rotate-180">
                                        <CaretDownOutlined/>
                                    </span>
                                </summary>
                                <ul className="mt-2 space-y-1 px-4">
                                    <li>
                                        <Link to={"/admin/link"}
                                              className={`transition flex items-center gap-2 rounded-lg px-4 py-2 text-sm ${SelectColorForDesk("/admin/link")}`}>
                                            管理友链
                                        </Link>
                                    </li>
                                    <li>
                                        <Link to={"/admin/link/add"}
                                              className={`transition flex items-center gap-2 rounded-lg px-4 py-2 text-sm ${SelectColorForDesk("/admin/link/add")}`}>
                                            添加友链
                                        </Link>
                                    </li>
                                    <li>
                                        <Link to={"/admin/link/location"}
                                              className={`transition flex items-center gap-2 rounded-lg px-4 py-2 text-sm ${SelectColorForDesk("/admin/link/location")}`}>
                                            位置管理
                                        </Link>
                                    </li>
                                    <li>
                                        <Link to={"/admin/link/color"}
                                              className={`transition flex items-center gap-2 rounded-lg px-4 py-2 text-sm ${SelectColorForDesk("/admin/link/color")}`}>
                                             颜色管理
                                        </Link>
                                    </li>
                                </ul>
                            </details>
                        </li>
                    </ul>
                </div>
                <div className="sticky inset-x-0 bottom-0 border-t border-gray-100">
                    <div className="transition flex items-center gap-2 bg-white p-4 hover:bg-gray-50">
                        <img alt="" className="size-10 rounded-full object-cover"
                             src="https://cravatar.cn/avatar/76602d1259d6a5e0796933f5d0ff9b84?s=400"/>
                        <div className="text-xs">
                            <strong className="block font-medium">{userCurrent.username}</strong>
                            <span>{userCurrent.email}</span>
                        </div>
                    </div>
                </div>
            </div>
            {/* 手机菜单样式 */}
            <div className={"fixed md:hidden w-lvw bottom-0 left-0 p-3 bg-white"}>
                <div className={"flex items-center justify-center"}>
                    <nav className="grid grid-cols-4 gap-3 w-full" aria-label="Tabs">
                        <Link to={"/admin/dashboard"}
                              className={`transition shrink-0 rounded-lg flex items-center justify-center gap-1 ${SelectColorForPhone("dashboard")}`}>
                            <HomeOutlined/>
                            <span>首页</span>
                        </Link>
                        <Link to={"/admin/link"}
                              className={`transition shrink-0 rounded-lg flex items-center justify-center gap-1 ${SelectColorForPhone("link")}`}>
                            <LinkOutlined/>
                            <span>友链</span>
                        </Link>
                        <Link to={"/admin/other"}
                              className={`transition shrink-0 rounded-lg flex items-center justify-center gap-1 ${SelectColorForPhone("other")}`}>
                            <BarsOutlined/>
                            <span>其他</span>
                        </Link>
                        <Link to={"/admin/setting"}
                              className={`transition shrink-0 rounded-lg flex items-center justify-center gap-1 ${SelectColorForPhone("setting")}`}
                              aria-current="page">
                            <SettingOutlined/>
                            <span>其他</span>
                        </Link>
                    </nav>
                </div>
            </div>
        </>
    );
}
