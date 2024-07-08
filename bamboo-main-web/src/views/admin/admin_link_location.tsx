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
import {
    AdminAddLocationAPI,
    AdminDelLocationAPI,
    AdminEditLocationAPI,
    AdminGetLocationAPI
} from "../../resources/ts/apis/api_link.ts";
import {message, Modal} from "antd";
import {LocationDO, LocationGetAdminEntity} from "../../resources/ts/models/entity/location_get_admin_entity.ts";
import {LocationDeleteDTO} from "../../resources/ts/models/dto/location_delete.ts";
import {CheckCircleOutlined, CloseCircleOutlined} from "@ant-design/icons";
import {LocationAddDTO} from "../../resources/ts/models/dto/location_add.ts";
import {Util} from "../../resources/utils/process_util.ts";

export function AdminLinkLocation() {
    const [locationInfo, setLocationInfo] = useState({} as LocationGetAdminEntity);

    const [isEditModalOpen, setIsEditModalOpen] = useState(false);
    const [isDelModalOpen, setIsDelModalOpen] = useState(false);
    const [isAddModalOpen, setIsAddModalOpen] = useState(false);

    const [editLocation, setEditLocation] = useState({} as LocationDO);
    const [delLocation, setDelLocation] = useState({} as LocationDeleteDTO);
    const [addLocation, setAddLocation] = useState({} as LocationAddDTO);

    const [locationWeb, setLocationWeb] = useState([] as JSX.Element[]);
    const [linkLocation, setLinkLocation] = useState([] as JSX.Element[]);


    useEffect(() => {
        setTimeout(async () => {
            const getRes = await AdminGetLocationAPI();
            if (getRes?.output === "Success") {
                setLocationInfo(getRes.data!);
            }
        });
        setTimeout(async () => {
            const getRes = await AdminGetLocationAPI();
            if (getRes?.output === "Success") {
                const locationList: JSX.Element[] = [];
                for (let i = 0; i < getRes.data!.locations.length; i++) {
                    locationList.push(<option key={i} value={getRes.data!.locations[i].id}>
                        {getRes.data!.locations[i].display_name}
                    </option>);
                }
                setLinkLocation(locationList);
            } else {
                message.warning(getRes?.error_message);
            }
        });
    }, []);

    useEffect(() => {
        if (!locationInfo.locations) {
            return;
        }
        setTimeout(async () => {
            // 写入信息进入 state
            const webData: JSX.Element[] = [];
            locationInfo.locations.forEach((resource) => {
                webData.push(<tr className="odd:bg-gray-50" key={resource.id}>
                    <td className="whitespace-nowrap px-4 py-2 font-medium text-gray-900">{resource.name}</td>
                    <td className="whitespace-nowrap px-4 py-2 text-gray-700">{resource.display_name}</td>
                    <td className="whitespace-nowrap px-4 py-2 text-gray-700 hidden md:table-cell">{whetherToDisplay(resource.reveal)}</td>
                    <td className="whitespace-nowrap px-4 py-2 text-gray-700 hidden xl:table-cell">{resource.created_at}</td>
                    <td className="whitespace-nowrap px-4 py-2 text-gray-700 hidden lg:table-cell">{resource.updated_at}</td>
                    <td className="whitespace-nowrap px-4 py-2 text-gray-700 flex gap-1 justify-end">
                        <button className={"rounded-md text-white px-2 bg-emerald-500 hover:bg-emerald-600"}
                                onClick={() => {
                                    setIsEditModalOpen(true);
                                    setEditLocation(resource);
                                }}>
                            编辑
                        </button>
                        <button className={"rounded-md text-white px-2 bg-red-500 hover:bg-red-600"}
                                onClick={() => {
                                    setIsDelModalOpen(true);
                                    setEditLocation(resource);
                                    setDelLocation({id: resource.id, move_id: 0});
                                }}>删除
                        </button>
                    </td>
                </tr>);
            });
            setLocationWeb(webData);
        });
    }, [locationInfo]);

    async function editSubmit() {
        // 提交数据进行修改操作
        const getRes = await AdminEditLocationAPI(editLocation);
        if (getRes?.output === "Success") {
            message.success("操作成功");
            // 关闭模态框
            setIsEditModalOpen(false);
            setEditLocation({} as LocationDO);
            // 更新数据
            const getRes = await AdminGetLocationAPI();
            if (getRes?.output === "Success") {
                setLocationInfo(getRes.data!);
            }
        } else {
            message.warning(getRes?.error_message);
        }
    }

    async function deleteSubmit() {
        const getRes = await AdminDelLocationAPI(delLocation);
        if (getRes?.output === "Success") {
            message.success("操作成功");
            // 关闭模态框
            setIsDelModalOpen(false);
            setDelLocation({} as LocationDeleteDTO);
            // 更新数据
            const getRes = await AdminGetLocationAPI();
            if (getRes?.output === "Success") {
                setLocationInfo(getRes.data!);
            }
        } else {
            message.warning(getRes?.error_message);
        }
    }

    async function addSubmit() {
        const getRes = await AdminAddLocationAPI(addLocation);
        if (getRes?.output === "Success") {
            message.success("操作成功");
            // 关闭模态框
            setIsAddModalOpen(false);
            setAddLocation({} as LocationAddDTO);
            // 更新数据
            const getRes = await AdminGetLocationAPI();
            if (getRes?.output === "Success") {
                setLocationInfo(getRes.data!);
            }
        } else {
            message.warning(getRes?.error_message);
        }
    }

    document.title = "竹叶 - 位置管理"

    return (<div className={"grid grid-cols-12 gap-3"}>
            <div className={"col-span-full text-xl font-bold mb-3"}>位置管理</div>
            <div className={"col-span-full flex justify-end"}>
                <button onClick={() => setIsAddModalOpen(true)}
                        className={"transition col-span-full rounded-lg text-white bg-emerald-500 hover:bg-emerald-600 px-4 py-1"}>
                    添加位置
                </button>
            </div>
            <div className={"col-span-full"}>
                <div className="overflow-x-auto rounded-lg shadow-md">
                    <table className="min-w-full divide-y-2 divide-gray-200 bg-white text-sm">
                        <thead className={"text-left font"}>
                        <tr>
                            <th className="px-4 py-2 font-medium text-gray-900">名字</th>
                            <th className="px-4 py-2 font-medium text-gray-900">别名</th>
                            <th className="hidden md:table-cell px-4 py-2 font-medium text-gray-900">是否展示</th>
                            <th className="hidden xl:table-cell px-4 py-2 font-medium text-gray-900">创建时间</th>
                            <th className="hidden lg:table-cell px-4 py-2 font-medium text-gray-900">修改时间</th>
                            <th className="px-4 py-2 font-medium text-gray-900 text-end">操作</th>
                        </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-200">
                        {locationWeb}
                        </tbody>
                    </table>
                </div>
            </div>
            <Modal title="位置添加" open={isAddModalOpen} onCancel={() => setIsAddModalOpen(false)}
                   footer={<div className={"flex gap-1 lg:gap-3 justify-end"}>
                       <button
                           className={"rounded-lg transition px-4 py-1 text-white bg-red-400 hover:bg-red-500"}
                           onClick={() => {
                               setIsAddModalOpen(false);
                               setAddLocation({} as LocationAddDTO);
                           }}>
                           取消
                       </button>
                       <button
                           className={"rounded-lg transition px-4 py-1 text-white bg-emerald-500 hover:bg-emerald-600"}
                           onClick={addSubmit}>
                           添加
                       </button>
                   </div>}>
                <div className={"grid grid-cols-2 gap-3 py-3"}>
                    <div className={"col-span-full sm:col-span-1"}>
                        <label htmlFor="LocationName"
                               className="block text-xs font-medium text-gray-700">位置识别名称</label>
                        <input
                            type="text"
                            id="LocationName"
                            name="name"
                            placeholder={"our"}
                            onChange={data => {
                                addLocation.name = data.target.value;
                            }}
                            className="mt-1 w-full rounded-md border-gray-200 shadow-sm text-sm"
                        />
                    </div>
                    <div className={"col-span-full sm:col-span-1"}>
                        <label htmlFor="LocationDisplayName"
                               className="block text-xs font-medium text-gray-700">位置展示名称</label>
                        <input
                            type="text"
                            id="LocationDisplayName"
                            name="display_name"
                            placeholder={"我自己的"}
                            onChange={data => {
                                addLocation.display_name = data.target.value;
                            }}
                            className="mt-1 w-full rounded-md border-gray-200 shadow-sm text-sm"
                        />
                    </div>
                    <div className={"col-span-full"}>
                        <label htmlFor="LocationSort" className="block text-xs font-medium text-gray-700">优先级</label>
                        <input
                            type="number"
                            id="LocationSort"
                            name="sort"
                            placeholder={"1"}
                            onChange={data => {
                                addLocation.sort = Number(data.target.value);
                            }}
                            className="mt-1 w-full rounded-md border-gray-200 shadow-sm text-sm"
                        />
                    </div>
                    <div className={"col-span-full"}>
                        <label htmlFor="LocationDesc"
                               className="block text-sm font-medium text-gray-700">位置描述</label>
                        <textarea
                            id="LocationDesc"
                            name="description"
                            rows={2}
                            placeholder={"这是我自己的展示地方，用于放我自己提供的服务"}
                            onChange={data => {
                                addLocation.description = data.target.value;
                            }}
                            className="mt-2 w-full rounded-lg border-gray-200 align-top shadow-sm sm:text-sm"
                        />
                    </div>
                    <div className={"col-span-full"}>
                        <label htmlFor="Option1" className="flex cursor-pointer items-center justify-end gap-3">
                            <strong className="font-medium text-gray-900">是否对外可见</strong>
                            <input
                                type="checkbox"
                                name="reveal"
                                defaultChecked={true}
                                onChange={data => {
                                    addLocation.reveal = Boolean(data.target.value);
                                }}
                                className={"transition size-4 rounded border-gray-300 text-emerald-500"}
                            />
                        </label>
                    </div>
                </div>
            </Modal>
            <Modal title="位置编辑" open={isEditModalOpen} onCancel={() => setIsEditModalOpen(false)}
                   footer={<div className={"flex gap-1 lg:gap-3 justify-end"}>
                       <button
                           className={"rounded-lg transition px-4 py-1 text-white bg-emerald-500 hover:bg-emerald-600"}
                           onClick={() => {
                               setIsEditModalOpen(false);
                               setEditLocation({} as LocationDO);
                           }}>
                           取消
                       </button>
                       <button
                           className={"rounded-lg transition px-4 py-1 text-white bg-red-400 hover:bg-red-500"}
                           onClick={editSubmit}>
                           修改
                       </button>
                   </div>}>
                <div className={"grid grid-cols-2 gap-3 py-3"}>
                    <div className={"col-span-full sm:col-span-1"}>
                        <label htmlFor="id" className="block text-xs font-medium text-gray-700">序列号</label>
                        <input
                            type="number"
                            id="id"
                            value={editLocation.id || ""}
                            className="mt-1 w-full rounded-md border-gray-200 shadow-sm text-sm bg-gray-100"
                            readOnly={true}
                        />
                    </div>
                    <div className={"col-span-full sm:col-span-1"}>
                        <label htmlFor="created_at" className="block text-xs font-medium text-gray-700">创建时间</label>
                        <input
                            type="text"
                            id="created_at"
                            value={editLocation.created_at || ""}
                            className="mt-1 w-full rounded-md border-gray-200 shadow-sm text-sm bg-gray-100"
                            readOnly={true}
                        />
                    </div>
                    <div className={"col-span-full sm:col-span-1"}>
                        <label htmlFor="name"
                               className="block text-xs font-medium text-gray-700">位置识别名称</label>
                        <input
                            type="text"
                            id="name"
                            value={editLocation.name || ""}
                            onChange={(event) => setEditLocation(Util.handleInputChange(editLocation, event))}
                            className="mt-1 w-full rounded-md border-gray-200 shadow-sm text-sm"
                        />
                    </div>
                    <div className={"col-span-full sm:col-span-1"}>
                        <label htmlFor="display_name"
                               className="block text-xs font-medium text-gray-700">位置展示名称</label>
                        <input
                            type="text"
                            id="display_name"
                            value={editLocation.display_name || ""}
                            onChange={(event) => setEditLocation(Util.handleInputChange(editLocation, event))}
                            className="mt-1 w-full rounded-md border-gray-200 shadow-sm text-sm"
                        />
                    </div>
                    <div className={"col-span-full"}>
                        <label htmlFor="sort" className="block text-xs font-medium text-gray-700">优先级</label>
                        <input
                            type="number"
                            id="sort"
                            value={editLocation.sort || ""}
                            onChange={(event) => setEditLocation(Util.handleInputChange(editLocation, event))}
                            className="mt-1 w-full rounded-md border-gray-200 shadow-sm text-sm"
                        />
                    </div>
                    <div className={"col-span-full"}>
                        <label htmlFor="description"
                               className="block text-sm font-medium text-gray-700">位置描述</label>
                        <textarea
                            id="description"
                            rows={2}
                            value={editLocation.description || ""}
                            onChange={(event) => setEditLocation(Util.handleInputChange(editLocation, event))}
                            className="mt-2 w-full rounded-lg border-gray-200 align-top shadow-sm sm:text-sm"
                        />
                    </div>
                    <div className={"col-span-full"}>
                        <label htmlFor="reveal" className="flex cursor-pointer items-center justify-end gap-3">
                            <strong className="font-medium text-gray-900">是否对外可见</strong>
                            <input
                                type="checkbox"
                                id={"reveal"}
                                checked={editLocation.reveal || false}
                                onChange={(event) => setEditLocation(Util.handleInputChange(editLocation, event))}
                                className={"transition size-4 rounded border-gray-300 text-emerald-500"}
                            />
                        </label>
                    </div>
                </div>
            </Modal>
            <Modal title="位置删除" open={isDelModalOpen} onCancel={() => setIsDelModalOpen(false)}
                   footer={<div className={"flex gap-1 lg:gap-3 justify-end"}>
                       <button
                           className={"rounded-lg transition px-4 py-1 text-white bg-emerald-500 hover:bg-emerald-600"}
                           onClick={() => {
                               setIsDelModalOpen(false);
                               setEditLocation({} as LocationDO);
                           }}>
                           取消
                       </button>
                       <button
                           className={"rounded-lg transition px-4 py-1 text-white bg-red-400 hover:bg-red-500"}
                           onClick={deleteSubmit}>
                           删除
                       </button>
                   </div>}>
                <div className={"grid gap-1 py-3"}>
                    <div className={"text-gray-500"}>您真的要删除 <span
                        className={"font-bold text-black"}>{editLocation.name}【{editLocation.display_name}】</span> 吗？
                    </div>
                    <div className={"font-light text-red-500"}>这将会失去很久很久......</div>
                    <div>
                        <label htmlFor="HeadlineAct"
                               className="block text-sm font-medium text-gray-900">该模块删除后，默认链接将移到</label>
                        <select
                            name="HeadlineAct"
                            id="HeadlineAct"
                            onChange={(data) => {
                                setDelLocation({id: delLocation.id, move_id: Number(data.target.value)});
                            }}
                            className="mt-1.5 w-full rounded-lg border-gray-300 text-gray-700 sm:text-sm"
                        >
                            <option value="">默认(优先级高一个)</option>
                            {linkLocation}
                        </select>
                    </div>
                </div>
            </Modal>
        </div>);
}

function whetherToDisplay(display: boolean) {
    if (display) {
        return (<span
                className="inline-flex gap-1 items-center justify-center rounded-md bg-emerald-100 px-2 text-emerald-700"
            >
              <CheckCircleOutlined/>
              <p className="whitespace-nowrap text-sm">展示</p>
            </span>);
    } else {
        return (<span
                className="inline-flex gap-1 items-center justify-center rounded-md bg-red-100 px-2 text-red-700"
            >
              <CloseCircleOutlined/>
              <p className="whitespace-nowrap text-sm">关闭</p>
            </span>);
    }
}
