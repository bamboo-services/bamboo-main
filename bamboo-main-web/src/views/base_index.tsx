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

import BackgroundWithIndex from "../resources/ts/body_background";
import {Link} from "react-router-dom";
import myAvatar from "../resources/images/my_avatar.png";
import {useEffect, useState} from "react";
import {InfoAPI} from "../resources/ts/apis/api_info.ts";
import {message} from "antd";

export default function BaseIndex() {

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

    document.title = `${systemInfo.site.site_name} - ${systemInfo.blogger.nick}`;

    return (
        <div style={BackgroundWithIndex} className={"h-dvh w-lvw grid"}>
            <div className={"grid grid-cols-4 gap-8 justify-center items-center px-8 lg:px-32"}>
                <div className={"hidden lg:flex col-span-1 items-center"}>
                    <img alt="UserAvatar" className={"rounded-full size-auto object-cover lg:h-48 xl:h-64"}
                         src={myAvatar} draggable={false}/>
                </div>
                <div className={"col-span-4 lg:col-span-3"}>
                    <div className={"text-center grid gap-3"}>
                        <h1 className="bg-gradient-to-r from-green-300 via-blue-500 to-purple-600 bg-clip-text text-4xl font-extrabold text-transparent sm:text-5xl"
                            style={{textShadow: "1px 1px 4px rgba(38,164,192,0.32)"}}>
                            {systemInfo.site.site_name}
                        </h1>
                        <div className={"lg:hidden items-center flex justify-center"}>
                            <img alt="UserAvatar" className={"rounded-xl w-auto h-32"}
                                 src={myAvatar} draggable={false}/>
                        </div>
                        <p className={"lg:text-xl/relaxed"}>
                            {systemInfo.blogger.description}
                        </p>
                        <div className={"flex flex-wrap justify-center gap-4 pt-8 md:pt-6 lg:pt-4"}>
                            <Link
                                className={"transition block rounded bg-blue-500 px-12 py-3 text-sm font-medium text-white hover:bg-blue-600 focus:outline-none focus:ring sm:w-auto shadow-xl shadow-blue-500/50"}
                                to={"https://blog.x-lf.com"}
                            >
                                去我的博客吧
                            </Link>
                            <Link to={"/me/about"}
                                  className={"transition block rounded bg-green-500 px-12 py-3 text-sm font-medium text-white hover:bg-green-600 focus:outline-none focus:ring sm:w-auto shadow-xl shadow-green-500/55"}>
                                了解我的更多
                            </Link>
                        </div>
                    </div>
                </div>
            </div>
            <footer
                className={"hidden md:flex absolute inset-x-0 bottom-0 justify-between items-end p-3 text-gray-500"}>
                <div className={"grid"}>
                    <Link to={'/auth/login'}>账户登录</Link>
                    <span>Copyright (C) 2016-2024 筱锋xiao_lfeng. All Rights Reserved.</span>
                </div>
                <div className={"grid text-end"}>
                    <Link to={"https://beian.miit.gov.cn/#/Integrated/index"} target={"_blank"}>粤ICP备 2022014822
                        号</Link>
                    <Link to={"https://beian.mps.gov.cn/#/query/webSearch"} target={"_blank"}>粤公网安备
                        44030702003207 号</Link>
                </div>
            </footer>
            <footer className={"grid md:hidden text-gray-500 text-center pt-8 absolute inset-x-0 bottom-0 pb-3"}>
                <Link to={'/auth/login'}>账户登录</Link>
                <Link to={"https://beian.miit.gov.cn/#/Integrated/index"} target={"_blank"}>
                    粤ICP备 2022014822 号
                </Link>
            </footer>
        </div>
    );
}
