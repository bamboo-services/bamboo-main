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

import {JSX, useEffect, useState} from "react";
import {AdminGetLinkAPI} from "../../resources/ts/apis/api_link.ts";
import {message} from "antd";
import {LinkDO} from "../../resources/ts/models/entity/link_get_admin_entity.ts";
import {Link} from "react-router-dom";

export function AdminLink() {
    const [linkCount, setLinkCount] = useState(0);
    const [linkList, setLinkList] = useState({} as LinkDO[]);
    const [webLink, setWebLink] = useState([<>
        <div className={"text-center col-span-full text-xl font-bold"}>当前不存在链接</div>
    </>]);
    const [webLinkTable, setWebLinkTable] = useState([<>
        <tr className="odd:bg-gray-50">
            <td className="whitespace-nowrap px-4 py-2 font-medium text-gray-900">不存在链接</td>
            <td className="whitespace-nowrap px-4 py-2 text-gray-700">--</td>
            <td className="whitespace-nowrap px-4 py-2 text-gray-700 hidden md:block">--</td>
            <td className="whitespace-nowrap px-4 py-2 text-gray-700">--</td>
        </tr>
    </>]);

    useEffect(() => {
        setTimeout(async () => {
            const getRes = await AdminGetLinkAPI();
            if (getRes?.output === "Success") {
                // 数据的提取整理
                setLinkCount(getRes.data!.total);
                setLinkList(getRes.data!.links);
            } else {
                message.warning(getRes?.error_message);
            }
        });
    }, []);

    useEffect(() => {
        if (linkCount > 0) {
            setWebLink(linkList.map((link: LinkDO) => {
                return (
                    <div key={link.id}
                         className={"transition block bg-white rounded-lg p-4 shadow-sm shadow-indigo-100 hover:scale-105"}>
                        <div className={"grid gap-1"}>
                            <div className={"flex gap-3"}>
                                <div className={"h-12 w-12 overflow-hidden flex rounded-full"}>
                                    {selectImageIsDirectOrCDN(link)}
                                </div>
                                <div className={"w-full"}>
                                    <div className={"flex justify-between gap-1 items-center"}>
                                        <div className={"text-xl font-bold truncate"}>{link.site_name}</div>
                                        <Link to={`/admin/link/edit/${link.id}`}
                                              className={"block rounded-lg bg-red-400 px-2 py-1 text-sm font-medium text-white"}>
                                            编辑
                                        </Link>
                                    </div>
                                    <div className={"text-sm text-gray-400 truncate"}>{link.site_url}</div>
                                </div>
                            </div>
                            <div className={"text-sm text-gray-700 truncate"}>{link.site_description}</div>
                        </div>
                    </div>
                );
            }));
            setWebLinkTable(linkList.map((link: LinkDO) => {
                return (
                    <tr key={link.id} className="odd:bg-gray-50">
                        <td className="whitespace-nowrap px-4 py-2 font-medium text-gray-900">{link.site_name}</td>
                        <td className="whitespace-nowrap px-4 py-2 text-gray-700">{link.site_url}</td>
                        <td className="whitespace-nowrap px-4 py-2 text-gray-700 hidden md:block">{link.site_url}</td>
                        <td className="whitespace-nowrap px-4 py-2 text-gray-700">
                            <Link to={`/admin/link/edit/${link.id}`} className={"text-blue-500"}>编辑</Link>
                        </td>
                    </tr>
                );
            }));
        }
    }, [linkList]);

    document.title = "竹叶 - 友链管理";

    return (
        <div className={"grid grid-cols-12 gap-6"}>
            <div className={"col-span-12 text-xl font-bold"}>友链管理</div>
            <div className={"col-span-12 lg:hidden"}>
                <div
                    className="transition block bg-white rounded-lg p-4 shadow-sm shadow-indigo-100 hover:scale-105 text-center">
                    213
                </div>
            </div>
            <div className={"col-span-12 lg:col-span-8"}>
                <div className={"hidden lg:grid grid-cols-1 md:grid-cols-2 gap-3"}>
                    {webLink}
                </div>
                <div className="lg:hidden overflow-x-auto rounded-lg shadow shadow-indigo-100">
                    <table className="min-w-full divide-y-2 divide-gray-200 bg-white text-sm">
                        <thead className="text-left rtl:text-right">
                        <tr>
                            <th className="whitespace-nowrap px-4 py-2 font-medium text-gray-900">站点</th>
                            <th className="whitespace-nowrap px-4 py-2 font-medium text-gray-900 hidden md:block">站长</th>
                            <th className="whitespace-nowrap px-4 py-2 font-medium text-gray-900">地址</th>
                            <th className="whitespace-nowrap px-4 py-2 font-medium text-gray-900">操作</th>
                        </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-200">
                        {webLinkTable}
                        </tbody>
                    </table>
                </div>
            </div>
            <div className={"hidden lg:block col-span-4 gap-3"}>
                <div
                    className="transition block bg-white rounded-lg p-4 shadow-sm shadow-indigo-100">
                    <div className={"grid grid-cols-12 gap-3"}>
                        <div className={"col-span-12 text-lg font-medium"}>条件查询</div>
                        <div className={"col-span-12"}>
                            <div>
                                <select
                                    name="HeadlineAct"
                                    id="HeadlineAct"
                                    className="mt-1.5 w-full rounded-lg border-gray-300 text-gray-700 sm:text-sm"
                                >
                                    <option value="">Please select</option>
                                    <option value="JM">John Mayer</option>
                                    <option value="SRV">Stevie Ray Vaughn</option>
                                    <option value="JH">Jimi Hendrix</option>
                                    <option value="BBK">B.B King</option>
                                    <option value="AK">Albert King</option>
                                    <option value="BG">Buddy Guy</option>
                                    <option value="EC">Eric Clapton</option>
                                </select>
                            </div>
                        </div>
                        <div className={"col-span-12"}>
                            <input
                                type="email"
                                id="UserEmail"
                                placeholder="john@rhcp.com"
                                className="mt-1 w-full rounded-lg border-gray-200 shadow-sm sm:text-sm"
                            />
                        </div>
                        <div className={"col-span-12 text-lg font-medium"}>友链状态</div>
                        <div className={"col-span-12"}>
                            <div className="flow-root">
                                <dl className="-my-3 divide-y divide-gray-100 text-sm">
                                    <div className="grid grid-cols-1 gap-1 py-3 sm:grid-cols-3 sm:gap-4">
                                        <dt className="font-medium text-gray-900">Title</dt>
                                        <dd className="text-gray-700 sm:col-span-2">Mr</dd>
                                    </div>

                                    <div className="grid grid-cols-1 gap-1 py-3 sm:grid-cols-3 sm:gap-4">
                                        <dt className="font-medium text-gray-900">Name</dt>
                                        <dd className="text-gray-700 sm:col-span-2">John Frusciante</dd>
                                    </div>

                                    <div className="grid grid-cols-1 gap-1 py-3 sm:grid-cols-3 sm:gap-4">
                                        <dt className="font-medium text-gray-900">Occupation</dt>
                                        <dd className="text-gray-700 sm:col-span-2">Guitarist</dd>
                                    </div>

                                    <div className="grid grid-cols-1 gap-1 py-3 sm:grid-cols-3 sm:gap-4">
                                        <dt className="font-medium text-gray-900">Salary</dt>
                                        <dd className="text-gray-700 sm:col-span-2">$1,000,000+</dd>
                                    </div>

                                    <div className="grid grid-cols-1 gap-1 py-3 sm:grid-cols-3 sm:gap-4">
                                        <dt className="font-medium text-gray-900">Bio</dt>
                                        <dd className="text-gray-700 sm:col-span-2">
                                            Lorem ipsum dolor, sit amet consectetur adipisicing elit. Et facilis debitis
                                            explicabo
                                            doloremque impedit nesciunt dolorem facere, dolor quasi veritatis quia fugit
                                            aperiam
                                            aspernatur neque molestiae labore aliquam soluta architecto?
                                        </dd>
                                    </div>
                                </dl>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}

function selectImageIsDirectOrCDN(link: LinkDO): JSX.Element {
    if (link.cdn_logo_url !== "") {
        return (
            <img src={link.cdn_logo_url}
                 className={"items-center justify-center w-full h-full object-cover"} alt={`link_${link.id}`}/>
        );
    }
    return (
        <img src={link.site_logo}
             className={"items-center justify-center w-full h-full object-cover"} alt={`link_${link.id}`}/>
    );
}
