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

import {BrowserRouter, Route, Routes} from "react-router-dom";
import React from "react";
import BaseIndex from "./views/base_index.tsx";
import BaseAbout from "./views/base_about.tsx";
import BaseAuth from "./views/base_auth.tsx";
import AuthLogin from "./views/auth/auth_login.tsx"
import {BaseAdmin} from "./views/base_admin.tsx";
import {AdminDashboard} from "./views/admin/admin_dashboard.tsx";

const AppRoutes: React.FC = () => {
    return (
        <BrowserRouter>
            <Routes>
                <Route path={"/"} element={<BaseIndex/>}/>
                <Route path={"about"} element={<BaseAbout/>}/>
                <Route path={"auth"} element={<BaseAuth/>}>
                    <Route path={"login"} element={<AuthLogin/>}/>
                </Route>
                <Route path={"admin"} element={<BaseAdmin/>}>
                    <Route path={"dashboard"} element={<AdminDashboard/>}/>
                </Route>
            </Routes>
        </BrowserRouter>
    );
};

export default AppRoutes;