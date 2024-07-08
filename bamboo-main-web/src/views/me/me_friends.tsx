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

import {Link} from "react-router-dom";
import {JSX, useEffect, useState} from "react";
import {LinkGetEntity} from "../../resources/ts/models/entity/link_get_entity.ts";
import {GetLinkAPI} from "../../resources/ts/apis/api_link.ts";
import {message} from "antd";
import {LinkOutlined} from "@ant-design/icons";
import {Util} from "../../resources/utils/process_util.ts";

export function MeFriends() {
    const [getLink, setGetLink] = useState({} as LinkGetEntity);
    const [webReveal, setWebReveal] = useState([] as JSX.Element[]);

    useEffect(() => {
        setTimeout(async () => {
            const getRes = await GetLinkAPI();
            if (getRes?.output === "Success") {
                setGetLink(getRes.data!);
            } else {
                message.warning(getRes?.error_message);
            }
        }, 1);
    }, []);

    useEffect(() => {
        try {
            // 对数据进行循环获取操作
            const locationWeb: JSX.Element[] = [];
            for (const locationElement of getLink.locations) {
                const locationLinks: JSX.Element[] = [];
                // 对内链接数据进行循环遍历
                for (const linkElement of locationElement.links) {
                    locationLinks.push(
                        <Link to={linkElement.site_url} target={"_blank"} key={linkElement.id}
                              className={"transition rounded-lg bg-white shadow-md shadow-white hover:scale-105 grid gap-1 text-center justify-center p-3"}>
                            <div className={"flex justify-center"}>
                                <img src={Util.TwoUrlSelectNoEmpty(linkElement.site_logo, linkElement.cdn_logo_url)}
                                     alt={""} draggable={false}
                                     className={"rounded-full size-12 lg:size-16 shadow-lg shadow-gray-100"}/>
                            </div>
                            <div className={"text-lg font-bold truncate"}>{linkElement.site_name}</div>
                            <div
                                className={"text-sm font-light text-gray-400 truncate"}>{linkElement.site_description}</div>
                        </Link>
                    );
                }
                // 对链接内容进行插入
                locationWeb.push(
                    <div className={"grid gap-3"} key={locationElement.id}>
                        <div className={"grid"}>
                            <div className={"flex gap-1"}>
                                <LinkOutlined/>
                                <div className={"text-xl font-medium"}>{locationElement.display_name}</div>
                            </div>
                            <div className={"text-sm text-gray-500"}>{locationElement.description}</div>
                        </div>
                        <div className={"grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-3 mb-8"}>
                            {locationLinks}
                        </div>
                    </div>
                );
            }
            setWebReveal(locationWeb);
        } catch (err) {
            console.log();
        }
    }, [getLink]);

    document.title = "竹叶 - 友人帐"

    return (
        <div className={"grid md:rounded-lg md:bg-white md:bg-opacity-50 md:p-6 md:shadow-xl"}>
            {webReveal}
        </div>
    );
}
