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

import {Link, useNavigate, useParams} from "react-router-dom";
import {JSX, useEffect, useState} from "react";
import {
    AdminGetColorAPI,
    AdminGetLocationAPI,
    EditLinkAPI,
    GetSingleLinkAPI
} from "../../resources/ts/apis/api_link.ts";
import {message} from "antd";
import {LeftOutlined} from "@ant-design/icons";
import {Util} from "../../resources/utils/process_util.ts";
import {InnerLinkDTO} from "../../resources/ts/models/entity/link_get_entity.ts";
import {LocationGetAdminEntity} from "../../resources/ts/models/entity/location_get_admin_entity.ts";
import {ColorsEntity} from "../../resources/ts/models/entity/color_get_entity.ts";
import {SystemInfo} from "../../resources/ts/models/entity/info_get_entity.ts";

export function AdminLinkEdit(systemInfo: SystemInfo) {
    const getParams = useParams();
    const navigate = useNavigate();

    const [getLink, setGetLink] = useState({} as InnerLinkDTO);
    const [getLocation, setGetLocation] = useState({} as LocationGetAdminEntity);
    const [getColor, setGetColor] = useState({} as ColorsEntity);

    const [webLocationSelect, setWebLocationSelect] = useState([] as JSX.Element[]);
    const [webColorSelect, setWebColorSelect] = useState([] as JSX.Element[]);

    useEffect(() => {
        setTimeout(async () => {
            const getRes = await GetSingleLinkAPI(Number(getParams.id));
            if (getRes?.output === "Success") {
                // 数据获取成功
                setGetLink(getRes.data!);
            } else {
                message.warning("链接不存在");
                navigate("/admin/link");
            }
        });
        setTimeout(async () => {
            // 获取所有的位置信息
            const getRes = await AdminGetLocationAPI();
            if (getRes?.output === "Success") {
                // 数据获取成功
                setGetLocation(getRes.data!);
            } else {
                message.warning(getRes?.error_message);
            }
        });
        setTimeout(async () => {
            const getRes = await AdminGetColorAPI();
            if (getRes?.output === "Success") {
                // 数据获取成功
                setGetColor(getRes.data!);
            } else {
                message.warning(getRes?.error_message);
            }
        });
    }, [getParams, navigate]);

    useEffect(() => {
        try {
            const webLocationSelect: JSX.Element[] = [];
            getLocation.locations.forEach((data) => {
                webLocationSelect.push(
                    <option value={data.id} key={data.id}>{data.display_name} [{data.name}]</option>
                );
            });
            setWebLocationSelect(webLocationSelect);
        } catch (err) {
            console.log();
        }
    }, [getLocation]);

    useEffect(() => {
        try {
            const webColorSelect: JSX.Element[] = [];
            getColor.colors.forEach((data) => {
                webColorSelect.push(
                    <option value={data.id} key={data.id}>
                        {data.display_name} [{data.name}]
                    </option>
                );
            });
            setWebColorSelect(webColorSelect);
        } catch (err) {
            console.log();
        }
    }, [getColor]);

    async function editSubmit() {
        const getRes = await EditLinkAPI(getLink);
        if (getRes?.output === "Success") {
            message.success("修改成功");
            navigate("/admin/link");
        } else {
            message.warning(getRes?.error_message);
        }
    }

    document.title = `${systemInfo.info.site.site_name} - 编辑友链`;

    return (
        <div className={"grid grid-cols-12 gap-6"}>
            <div className={"col-span-full text-xl font-bold"}>友链管理 - 编辑</div>
            <div className={"col-span-full flex lg:hidden justify-between"}>
                <Link to={"/admin/link"}
                      className={"transition hover:scale-125"}>
                    <LeftOutlined/>
                </Link>
                <div onClick={editSubmit}
                     className={"rounded-lg transition bg-emerald-500 hover:bg-emerald-600 text-white text-center px-4 py-1"}>
                    确认修改
                </div>
            </div>
            <div className={"col-span-full lg:col-span-8"} key={"first"}>
                <div className={"bg-white shadow-sm shadow-indigo-100 rounded-lg grid grid-cols-4 gap-3 p-3"}>
                    <h1 className={"col-span-full text-lg lg:text-xl font-bold"}>基本信息</h1>
                    <div className="col-span-full flow-root">
                        <dl className="-my-3 divide-y divide-gray-100 text-sm">
                            <div
                                className="grid grid-cols-2 gap-1 py-2 even:bg-gray-50 sm:grid-cols-3 lg:grid-cols-4 sm:gap-4">
                                <dt className="font-medium text-gray-900">序列号</dt>
                                <dd className="text-gray-700 sm:col-span-2 lg:col-span-3">{getLink.id}</dd>
                            </div>
                            <div
                                className="grid grid-cols-2 gap-1 py-2 even:bg-gray-50 sm:grid-cols-3 lg:grid-cols-4 sm:gap-4">
                                <dt className="font-medium text-gray-900">期望位置</dt>
                                <dd className="text-gray-700 sm:col-span-2 lg:col-span-3">{getLink.desired_location}</dd>
                            </div>
                            <div
                                className="grid grid-cols-2 gap-1 py-2 even:bg-gray-50 sm:grid-cols-3 lg:grid-cols-4 sm:gap-4">
                                <dt className="font-medium text-gray-900">期望颜色</dt>
                                <dd className="text-gray-700 sm:col-span-2 lg:col-span-3">{getLink.desired_color}</dd>
                            </div>
                            <div
                                className="grid grid-cols-2 gap-1 py-2 even:bg-gray-50 sm:grid-cols-3 lg:grid-cols-4 sm:gap-4">
                                <dt className="font-medium text-gray-900">站长留言</dt>
                                <dd className="text-gray-700 sm:col-span-2 lg:col-span-3">{getLink.webmaster_remark}</dd>
                            </div>
                            <div
                                className="grid grid-cols-2 gap-1 py-2 even:bg-gray-50 sm:grid-cols-3 lg:grid-cols-4 sm:gap-4">
                                <dt className="font-medium text-gray-900">状态</dt>
                                <dd className="text-gray-700 sm:col-span-2 lg:col-span-3">{getLink.status}</dd>
                            </div>
                            <div
                                className="grid grid-cols-2 gap-1 py-2 even:bg-gray-50 sm:grid-cols-3 lg:grid-cols-4 sm:gap-4">
                                <dt className="font-medium text-gray-900">创建时间</dt>
                                <dd className="text-gray-700 sm:col-span-2 lg:col-span-3">{getLink.created_at}</dd>
                            </div>
                            <div
                                className="grid grid-cols-2 gap-1 py-2 even:bg-gray-50 sm:grid-cols-3 lg:grid-cols-4 sm:gap-4">
                                <dt className="font-medium text-gray-900">修改时间</dt>
                                <dd className="text-gray-700 sm:col-span-2 lg:col-span-3">{getLink.updated_at}</dd>
                            </div>
                        </dl>
                    </div>
                    <h1 className={"col-span-full text-lg lg:text-xl font-bold"}>可修改信息</h1>
                    <span className="relative flex justify-center col-span-full">
                      <div
                          className="absolute inset-x-0 top-1/2 h-px -translate-y-1/2 bg-transparent bg-gradient-to-r from-transparent via-gray-500 to-transparent opacity-75"
                      />
                      <span className="relative z-10 bg-white px-6">站点配置</span>
                    </span>
                    <div className={"col-span-full lg:col-span-2"}>
                        <label htmlFor="webmaster_email" className="text-xs font-medium text-gray-700 flex gap-1">
                            <span>站长邮箱</span>
                        </label>
                        <input
                            type="text"
                            id="webmaster_email"
                            value={getLink.webmaster_email || ""}
                            onChange={(event) => setGetLink(Util.handleInputChange(getLink, event))}
                            className="mt-1 w-full rounded-md border-gray-200 shadow-sm text-sm"
                        />
                    </div>
                    <div className={"col-span-full lg:col-span-2"}>
                        <label htmlFor="service_provider" className="text-xs font-medium text-gray-700 flex gap-1">
                            <span>站点服务提供商</span>
                        </label>
                        <input
                            type="text"
                            id="service_provider"
                            value={getLink.service_provider || ""}
                            onChange={(event) => setGetLink(Util.handleInputChange(getLink, event))}
                            className="mt-1 w-full rounded-md border-gray-200 shadow-sm text-sm"
                        />
                    </div>
                    <span className="relative flex justify-center col-span-full">
                      <div
                          className="absolute inset-x-0 top-1/2 h-px -translate-y-1/2 bg-transparent bg-gradient-to-r from-transparent via-gray-500 to-transparent opacity-75"
                      />
                      <span className="relative z-10 bg-white px-6">站点信息</span>
                    </span>
                    <div className={"col-span-full lg:col-span-2"}>
                        <label htmlFor="site_name" className="text-xs font-medium text-gray-700 flex gap-1">
                            <span>站点名字</span>
                            <span className={"text-red-500"}>*</span>
                        </label>
                        <input
                            type="text"
                            id="site_name"
                            value={getLink.site_name || ""}
                            onChange={(event) => setGetLink(Util.handleInputChange(getLink, event))}
                            className="mt-1 w-full rounded-md border-gray-200 shadow-sm text-sm"
                        />
                    </div>
                    <div className={"col-span-full lg:col-span-2"}>
                        <label htmlFor="site_name" className="text-xs font-medium text-gray-700 flex gap-1">
                            <span>站点地址</span>
                            <span className={"text-red-500"}>*</span>
                        </label>
                        <input
                            type="text"
                            id="site_url"
                            value={getLink.site_url || ""}
                            onChange={(event) => setGetLink(Util.handleInputChange(getLink, event))}
                            className="mt-1 w-full rounded-md border-gray-200 shadow-sm text-sm"
                        />
                    </div>
                    <div className={"col-span-full lg:col-span-2"}>
                        <label htmlFor="site_logo" className="text-xs font-medium text-gray-700 flex gap-1">
                            <span>站点图标</span>
                            <span className={"text-red-500"}>*</span>
                        </label>
                        <input
                            type="text"
                            id="site_logo"
                            value={getLink.site_logo || ""}
                            onChange={(event) => setGetLink(Util.handleInputChange(getLink, event))}
                            className="mt-1 w-full rounded-md border-gray-200 shadow-sm text-sm"
                        />
                    </div>
                    <div className={"col-span-full lg:col-span-2"}>
                        <label htmlFor="cdn_logo_url" className="text-xs font-medium text-gray-700 flex gap-1">
                            <span>站点加速图标</span>
                        </label>
                        <input
                            type="text"
                            id="cdn_logo_url"
                            value={getLink.cdn_logo_url || ""}
                            onChange={(event) => setGetLink(Util.handleInputChange(getLink, event))}
                            className="mt-1 w-full rounded-md border-gray-200 shadow-sm text-sm"
                        />
                    </div>
                    <div className={"col-span-full"}>
                        <label htmlFor="site_rss_url" className="text-xs font-medium text-gray-700 flex gap-1">
                            <span>站点订阅地址</span>
                        </label>
                        <input
                            type="text"
                            id="site_rss_url"
                            value={getLink.site_rss_url || ""}
                            onChange={(event) => setGetLink(Util.handleInputChange(getLink, event))}
                            className="mt-1 w-full rounded-md border-gray-200 shadow-sm text-sm"
                        />
                    </div>
                    <div className={"col-span-full"}>
                        <label htmlFor="site_description"
                               className="text-xs font-medium text-gray-700 flex gap-1">
                            <span>站点描述</span>
                        </label>
                        <textarea
                            id="site_description"
                            rows={2}
                            value={getLink.site_description || ""}
                            onChange={(event) => setGetLink(Util.handleInputChange(getLink, event))}
                            className="mt-2 w-full rounded-lg border-gray-200 align-top shadow-sm sm:text-sm"
                        />
                    </div>
                    <span className="relative flex justify-center col-span-full">
                      <div
                          className="absolute inset-x-0 top-1/2 h-px -translate-y-1/2 bg-transparent bg-gradient-to-r from-transparent via-gray-500 to-transparent opacity-75"
                      />
                      <span className="relative z-10 bg-white px-6">位置管理</span>
                    </span>
                    <div className={"col-span-full lg:col-span-2"}>
                        <label htmlFor="location"
                               className="text-xs font-medium text-gray-700 flex gap-1">
                            <span>所选位置</span>
                            <span className={"text-red-500"}>*</span>
                        </label>
                        <select
                            name="HeadlineAct"
                            id="location"
                            value={getLink.location || ""}
                            onChange={(event) => setGetLink(Util.handleInputChange(getLink, event))}
                            className="mt-1.5 w-full rounded-lg border-gray-300 text-gray-700 sm:text-sm">
                            <option value="" key={"default"}>请选择位置</option>
                            {webLocationSelect}
                        </select>
                    </div>
                    <div className={"col-span-full lg:col-span-2"}>
                        <label htmlFor="color"
                               className="text-xs font-medium text-gray-700 flex gap-1">
                            <span>所选颜色</span>
                            <span className={"text-red-500"}>*</span>
                        </label>
                        <select
                            name="HeadlineAct"
                            id="color"
                            value={getLink.color || ""}
                            onChange={(event) => setGetLink(Util.handleInputChange(getLink, event))}
                            className="mt-1.5 w-full rounded-lg border-gray-300 text-gray-700 sm:text-sm">
                            <option value="" key={"default"}>请选择颜色</option>
                            {webColorSelect}
                        </select>
                    </div>
                    <div className={"col-span-full lg:col-span-2"}>
                        <label htmlFor="has_adv"
                               className="text-xs font-medium text-gray-700 flex gap-1">
                            <span>是否有广告</span>
                            <span className={"text-red-500"}>*</span>
                        </label>
                        <select
                            name="HeadlineAct"
                            id="has_adv"
                            value={String(getLink.has_adv) || ""}
                            onChange={(event) => setGetLink(Util.handleInputChange(getLink, event))}
                            className="mt-1.5 w-full rounded-lg border-gray-300 text-gray-700 sm:text-sm">
                            <option value="0">否</option>
                            <option value="1">是</option>
                        </select>
                    </div>
                    <span className="relative flex justify-center col-span-full">
                      <div
                          className="absolute inset-x-0 top-1/2 h-px -translate-y-1/2 bg-transparent bg-gradient-to-r from-transparent via-gray-500 to-transparent opacity-75"
                      />
                      <span className="relative z-10 bg-white px-6">其他信息</span>
                    </span>
                    <div className={"col-span-full"}>
                        <label htmlFor="remark"
                               className="text-xs font-medium text-gray-700 flex gap-1">
                            <span>我的备注</span>
                        </label>
                        <textarea
                            id="remark"
                            rows={2}
                            value={getLink.remark || ""}
                            onChange={(event) => setGetLink(Util.handleInputChange(getLink, event))}
                            className="mt-2 w-full rounded-lg border-gray-200 align-top shadow-sm sm:text-sm"
                        />
                    </div>
                </div>
            </div>
            <div className={"hidden lg:block col-span-4"} key={"second"}>
                <div className={"bg-white shadow-sm shadow-indigo-100 rounded-lg grid grid-cols-2 gap-3 p-3"}>
                    <div className={"col-span-full"}>
                        <div
                            className={"transition rounded-lg bg-white grid gap-1 text-center justify-center p-3 border"}>
                            <div className={"flex justify-center"}>
                                <img
                                    src={Util.TwoUrlSelectNoEmpty(getLink.site_logo, getLink.cdn_logo_url)}
                                    alt={""} draggable={false}
                                    className={"rounded-full size-12 lg:size-16 shadow-lg shadow-gray-100"}/>
                            </div>
                            <div className={"text-lg font-bold truncate"}>{getLink.site_name}</div>
                            <div className={"text-sm font-light text-gray-400 truncate"}>
                                {getLink.site_description}
                            </div>
                        </div>
                    </div>
                    <h1 className={"col-span-full text-lg lg:text-xl font-bold"}>操作</h1>
                    <Link to={"/admin/link"}
                          className={"rounded-lg transition bg-red-400 hover:bg-red-500 text-white text-center px-4 py-2"}>
                        取消修改
                    </Link>
                    <button
                        className={"rounded-lg transition bg-emerald-500 hover:bg-emerald-600 text-white text-center px-4 py-2"}
                        onClick={editSubmit}>
                        确认修改
                    </button>
                </div>
            </div>
        </div>
    );
}
