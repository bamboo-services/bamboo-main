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

import {AuthLoginDTO} from "../../resources/ts/models/dto/auth_login";
import {UserLoginAPI} from "../../resources/ts/apis/api_auth.ts";
import {message} from "antd";
import {useNavigate} from "react-router-dom";

export default function AuthLogin() {
    const loginForm = {remember: false} as AuthLoginDTO;
    const navigation = useNavigate();

    document.title = "竹叶 - 管理员登陆";

    /**
     * # 登录表单
     *
     * 用于提交用户登录的信息管理
     *
     * @param event 相应时间
     */
    async function onSubmit(event: { preventDefault: () => void; }) {
        event.preventDefault();
        // 检查数据是否为空
        if (!loginForm.user || !loginForm.pass) {
            message.warning("用户名或密码不能为空！");
            return;
        }
        // 提交登录信息
        const getRes = await UserLoginAPI(loginForm);
        if (getRes) {
            if (getRes.output === "Success") {
                message.success("登录成功！");
                // 保存 Token
                localStorage.setItem("UserToken", getRes.data?.token as string);
                // 路由跳转
                setTimeout(() => {
                    navigation("/admin/dashboard");
                }, 500);
            } else {
                message.info(getRes.error_message);
            }
        } else {
            message.error("登录失败！请检查后端连接是否正确！");
        }
    }

    return (
        <div className={"grid w-screen h-dvh items-center justify-center bg-gray-100"}>
            <div className={"grid gap-1 mx-3"}>
                <div className={"text-center"}>
                    <h1 className={"text-2xl font-bold text-emerald-600 sm:text-3xl"}>凌中的锋雨</h1>
                    <p className={"mx-auto mt-4 max-w-md text-gray-500"}>
                        一个人的价值，在于他贡献什么，而不是他能取得什么。不要渴望成为一个成功的人，而是应该努力做一个有价值的人。
                    </p>
                </div>
                <form onSubmit={onSubmit}
                      className={"mt-6 space-y-3 rounded-lg p-4 shadow-lg sm:p-6 lg:p-8 bg-white w-screen sm:max-w-screen-sm"}>
                    <p className={"text-center text-lg font-medium"}>用户登录</p>
                    <div>
                        <label htmlFor="user" className="sr-only">用户名</label>
                        <div className="relative">
                            <input
                                type="text"
                                className="w-full rounded-lg border-gray-200 p-3 text-sm shadow-sm"
                                placeholder="请输入用户名"
                                onChange={e => {
                                    loginForm.user = e.target.value;
                                }}
                            />
                        </div>
                    </div>
                    <div>
                        <label htmlFor="password" className="sr-only">Password</label>
                        <div className="relative">
                            <input
                                type="password"
                                className="w-full rounded-lg border-gray-200 p-3 text-sm shadow-sm"
                                placeholder="请输入密码"
                                onChange={e => {
                                    loginForm.pass = e.target.value;
                                }}
                            />
                        </div>
                    </div>
                    <div>
                        <label htmlFor="Option1" className="flex cursor-pointer items-start gap-4">
                            <div className="flex items-center">
                                &#8203;
                                <input
                                    type="checkbox"
                                    className="size-4 rounded border-emerald-300 text-emerald-500 ring ring-emerald-300 emer400 transition"
                                    onChange={e => {
                                        loginForm.remember = e.target.checked;
                                    }}
                                />
                            </div>
                            <div>
                                <strong className="font-medium text-gray-900">记住我 <span
                                    className={"text-gray-400"}>(7日免登录)</span></strong>
                            </div>
                        </label>
                    </div>
                    <div className={"w-full flex justify-center pt-4"}>
                        <button
                            type="submit"
                            className="block rounded-lg bg-emerald-500 hover:bg-emerald-600 px-12 py-2 text-sm font-medium text-white">
                            登录
                        </button>
                    </div>
                </form>
            </div>
        </div>

    );
}
