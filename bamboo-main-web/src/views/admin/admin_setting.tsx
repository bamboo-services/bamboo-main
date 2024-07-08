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

import {SystemInfo} from "../../resources/ts/models/entity/info_get_entity.ts";
import {MailSendAPI} from "../../resources/ts/apis/api_email.ts";
import {message} from "antd";
import {MailSendDTO} from "../../resources/ts/models/dto/mail_send.ts";
import {useEffect, useState} from "react";
import {AuthChangePasswordDTO} from "../../resources/ts/models/dto/auth_change_password.ts";
import {Util} from "../../resources/utils/process_util.ts";
import {UserChangePasswordAPI} from "../../resources/ts/apis/api_auth.ts";

export function AdminSetting(systemInfo: SystemInfo) {
    const [changePassword, setChangePassword] = useState({} as AuthChangePasswordDTO);
    const [mailCountdown, setMailCountdown] = useState(0);
    const [mailSend, setMailSend] = useState(false);

    useEffect(() => {
        if (localStorage.getItem("mailCountdown") !== "0" && mailCountdown === 0) {
            setMailCountdown(Number(localStorage.getItem("mailCountdown")));
            setMailSend(true);
        }
        if (mailSend) {
            // 设置禁止点击
            document.getElementById("send_email_submit")!.setAttribute("disabled", "true");
            document.getElementById("send_email_submit")!.classList.remove("bg-blue-500");
            document.getElementById("send_email_submit")!.classList.remove("hover:bg-blue-600");
            document.getElementById("send_email_submit")!.classList.add("bg-blue-800");
            // 倒计时
            document.getElementById("send_email_submit")!.innerText = `发送邮件(${mailCountdown})`;
            setTimeout(async () => {
                setMailCountdown(mailCountdown - 1);
                localStorage.setItem("mailCountdown", mailCountdown.toString());
                // 修改按钮文字
            }, 1000);
            if (mailCountdown <= 0) {
                localStorage.removeItem("mailCountdown");
                setMailSend(false);
                document.getElementById("send_email_submit")!.innerText = `发送邮件`;
                document.getElementById("send_email_submit")!.removeAttribute("disabled");
                document.getElementById("send_email_submit")!.classList.add("bg-blue-500");
                document.getElementById("send_email_submit")!.classList.add("hover:bg-blue-600");
                document.getElementById("send_email_submit")!.classList.remove("bg-blue-800");
            }
        }
    }, [mailCountdown, mailSend]);

    async function sendEmail() {
        const getRes = await MailSendAPI({
            to: changePassword.email,
            scene: "ChangePassword",
        } as MailSendDTO);
        if (getRes?.output === "Success") {
            message.success("邮件发送成功！");
            setMailSend(true);
            setMailCountdown(60);
        } else {
            message.warning(getRes?.error_message);
        }
    }

    async function resetPassword() {
        const getRes = await UserChangePasswordAPI(changePassword);
        if (getRes?.output === "Success") {
            message.success("密码修改成功");
            // eslint-disable-next-line @typescript-eslint/ban-ts-comment
            // @ts-expect-error
            document.getElementById("email")!.value = "";
            // eslint-disable-next-line @typescript-eslint/ban-ts-comment
            // @ts-expect-error
            document.getElementById("email_code")!.value = "";
            // eslint-disable-next-line @typescript-eslint/ban-ts-comment
            // @ts-expect-error
            document.getElementById("new_password")!.value = "";
        } else {
            message.warning(getRes?.error_message);
        }
    }

    document.title = `${systemInfo.info.site.site_name} - 系统设置`

    return (
        <div className={"grid grid-cols-12 gap-6"}>
            <div className={"col-span-full text-xl font-bold"}>系统设置</div>
            <div className={"col-span-full lg:col-span-6 space-y-3"}>
                <div className={"rounded-lg bg-white p-4 shadow-sm shadow-indigo-100 grid gap-3"}>
                    <div className={"flex justify-between"}>
                        <div className={"text-lg font-bold"}>系统设置</div>
                        <div className={"flex gap-3"}>
                            <button
                                className={"transition rounded-lg bg-emerald-500 text-white px-4 py-1 hover:bg-emerald-600"}>
                                保存设置
                            </button>
                        </div>
                    </div>
                </div>
            </div>
            <div className={"col-span-full lg:col-span-6 space-y-3"}>
                <div className={"rounded-lg bg-white p-4 shadow-sm shadow-indigo-100 grid gap-3"}>
                    <div className={"flex justify-between"}>
                        <div className={"text-lg font-bold"}>修改密码</div>
                        <div className={"flex gap-3"}>
                            <button onClick={sendEmail}
                                    id={"send_email_submit"}
                                    className={"transition rounded-lg bg-blue-500 text-white px-4 py-1 hover:bg-blue-600"}>
                                发送邮件
                            </button>
                            <button onClick={resetPassword}
                                    className={"transition rounded-lg bg-emerald-500 text-white px-4 py-1 hover:bg-emerald-600"}>
                                提交修改
                            </button>
                        </div>
                    </div>
                    <label
                        htmlFor="email"
                        className="relative block overflow-hidden border-b border-gray-200 bg-transparent pt-3 focus-within:border-blue-600">
                        <input type={"email"}
                               id={"email"}
                               onChange={(event) => setChangePassword(Util.handleInputChange(changePassword, event))}
                               className="peer h-8 w-full border-none bg-transparent p-0 placeholder-transparent focus:border-transparent focus:outline-none focus:ring-0 sm:text-sm"/>
                        <span
                            className="absolute start-0 top-2 -translate-y-1/2 text-xs text-gray-700 transition-all peer-placeholder-shown:top-1/2 peer-placeholder-shown:text-sm peer-focus:top-2 peer-focus:text-xs flex gap-1">
                            <span>邮箱</span>
                            <span className={"text-red-500"}>*</span>
                        </span>
                    </label>
                    <label
                        htmlFor="email_code"
                        className="relative block overflow-hidden border-b border-gray-200 bg-transparent pt-3 focus-within:border-blue-600">
                        <input
                            type="text"
                            id={"email_code"}
                            onChange={(event) => setChangePassword(Util.handleInputChange(changePassword, event))}
                            className="peer h-8 w-full border-none bg-transparent p-0 placeholder-transparent focus:border-transparent focus:outline-none focus:ring-0 sm:text-sm"/>
                        <span
                            className="absolute start-0 top-2 -translate-y-1/2 text-xs text-gray-700 transition-all peer-placeholder-shown:top-1/2 peer-placeholder-shown:text-sm peer-focus:top-2 peer-focus:text-xs flex gap-1">
                            <span>验证码</span>
                            <span className={"text-red-500"}>*</span>
                        </span>
                    </label>
                    <label
                        htmlFor="new_password"
                        className="relative block overflow-hidden border-b border-gray-200 bg-transparent pt-3 focus-within:border-blue-600">
                        <input
                            type="password"
                            id={"new_password"}
                            onChange={(event) => setChangePassword(Util.handleInputChange(changePassword, event))}
                            className="peer h-8 w-full border-none bg-transparent p-0 placeholder-transparent focus:border-transparent focus:outline-none focus:ring-0 sm:text-sm"/>
                        <span
                            className="absolute start-0 top-2 -translate-y-1/2 text-xs text-gray-700 transition-all peer-placeholder-shown:top-1/2 peer-placeholder-shown:text-sm peer-focus:top-2 peer-focus:text-xs flex gap-1">
                            <span>新密码</span>
                            <span className={"text-red-500"}>*</span>
                        </span>
                    </label>
                </div>
            </div>
        </div>
    );
}
