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

import {AdminDashboard} from "./admin/admin_dashboard.tsx";
import {Route, Routes, useNavigate} from "react-router-dom";
import {AdminSideMenuComponent} from "../components/admin/admin_side_menu.tsx";
import {useEffect, useState} from "react";
import {InfoAPI, InfoUserAPI} from "../resources/ts/apis/api_info.ts";
import {message} from "antd";
import {UserCurrentEntity} from "../resources/ts/models/entity/user_current_entity.ts";
import {AdminLink} from "./admin/admin_link.tsx";
import {AdminLinkEdit} from "./admin/admin_link_edit.tsx";
import {AdminLinkLocation} from "./admin/admin_link_location.tsx";
import {AdminSetting} from "./admin/admin_setting.tsx";
import {AdminLinkAdd} from "./admin/admin_link_add.tsx";

export function BaseAdmin() {
    const navigation = useNavigate();

    const [userCurrent, setUserInfo] = useState({} as UserCurrentEntity);
    const [systemInfo, setSystemInfo] = useState({
        site: {
            site_name: "竹叶",
            author: "筱锋xiao_lfeng",
            version: "",
            description: "",
            keywords: "",
        },
        blogger: {
            name: "",
            nick: "",
            email: "",
            description: "",
        }
    });

    useEffect(() => {
        if (localStorage.getItem("WebInfo") == null) {
            setTimeout(async () => {
                const getRes = await InfoAPI();
                if (getRes?.output === "Success") {
                    setSystemInfo(getRes.data!);
                    localStorage.setItem("WebInfo", JSON.stringify(getRes.data!));
                } else {
                    message.warning(getRes?.error_message);
                }
            });
        } else {
            setSystemInfo(JSON.parse(localStorage.getItem("WebInfo")!));
        }
    }, []);

    useEffect(() => {
        InfoUserAPI().then(async (getRes) => {
            if (getRes) {
                if (getRes.output !== "Success") {
                    message.warning(getRes?.error_message);
                    // 移除登录的内容
                    localStorage.removeItem("UserToken");
                    setTimeout(() => {
                        navigation("/auth/login");
                    }, 1000);
                } else {
                    setUserInfo(getRes.data!);
                }
            }
        })
    }, [navigation]);

    return (
        <div className={"grid grid-cols-12 bg-gray-100/75"}>
            <div className={"hidden md:block"}>
                <AdminSideMenuComponent username={userCurrent.username} email={userCurrent.email}
                                        uuid={userCurrent.uuid}/>
            </div>
            <div className={"md:ps-[200px] lg:ps-[250px] col-span-12 min-h-lvh bg-gray-100/75 pb-20 md:pb-0"}>
                <div className={"p-6"}>
                    <Routes>
                        <Route path={"dashboard"} element={<AdminDashboard info={systemInfo}/>}/>
                        <Route path={"link"} element={<AdminLink info={systemInfo}/>}/>
                        <Route path={"link/add"} element={<AdminLinkAdd info={systemInfo}/>}/>
                        <Route path={"link/edit/:id"} element={<AdminLinkEdit info={systemInfo}/>}/>
                        <Route path={"link/location"} element={<AdminLinkLocation info={systemInfo}/>}/>
                        <Route path={"setting"} element={<AdminSetting info={systemInfo}/>}/>
                    </Routes>
                </div>
            </div>
            <div className={"block md:hidden col-span-12"}>
                <AdminSideMenuComponent username={userCurrent.username} email={userCurrent.email}
                                        uuid={userCurrent.uuid}/>
            </div>
        </div>
    )
}
