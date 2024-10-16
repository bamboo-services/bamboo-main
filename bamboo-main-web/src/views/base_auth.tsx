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

import {Route, Routes, useLocation, useNavigate} from "react-router-dom";
import {useEffect, useState} from "react";
import AuthLogin from "./auth/auth_login";
import {InfoAPI, InfoUserAPI} from "../resources/ts/apis/api_info.ts";
import {message} from "antd";

export default function BaseAuth() {
    const navigate = useNavigate();
    const getLocation = useLocation();

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
        if (getLocation.pathname === "/auth") {
            navigate("/auth/login");
        }
        setTimeout(async () => {
            // 检查用户是否登录
            const getRes = await InfoUserAPI();
            if (getRes?.output === "Success") {
                message.success(`您好 ${getRes.data?.username} 用户`)
                setTimeout(() => {
                    navigate("/admin/dashboard");
                }, 500);
            }
        })
    });

    return (
        <Routes>
            <Route path={"login"} element={<AuthLogin info={systemInfo}/>}/>
        </Routes>
    );
}
