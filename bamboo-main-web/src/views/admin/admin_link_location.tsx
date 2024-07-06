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
import {AdminEditLocationAPI, AdminGetLocationAPI} from "../../resources/ts/apis/api_link.ts";
import {message, Modal} from "antd";
import {LocationDO, LocationGetAdminEntity} from "../../resources/ts/models/entity/location_get_admin_entity.ts";

export function AdminLinkLocation() {
    const [locationInfo, setLocationInfo] = useState({} as LocationGetAdminEntity);
    const [locationWeb, setLocationWeb] = useState([] as JSX.Element[]);
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [editLocation, setEditLocation] = useState({} as LocationDO);

    useEffect(() => {
        setTimeout(async () => {
            const getRes = await AdminGetLocationAPI();
            if (getRes?.output === "Success") {
                setLocationInfo(getRes.data!);
            }
        }, 1);
    }, []);

    useEffect(() => {
        if (!locationInfo.locations) {
            return;
        }
        setTimeout(async () => {
            // 写入信息进入 state
            const webData: JSX.Element[] = [];
            locationInfo.locations.forEach((resource) => {
                webData.push(
                    <tr className="odd:bg-gray-50" key={resource.id}>
                        <td className="whitespace-nowrap px-4 py-2 font-medium text-gray-900">{resource.name}</td>
                        <td className="whitespace-nowrap px-4 py-2 text-gray-700">{resource.display_name}</td>
                        <td className="whitespace-nowrap px-4 py-2 text-gray-700">{resource.reveal + ""}</td>
                        <td className="whitespace-nowrap px-4 py-2 text-gray-700">{resource.updated_at}</td>
                        <td className="whitespace-nowrap px-4 py-2 text-gray-700 flex gap-1 justify-end">
                            <button className={"rounded-md text-white px-2 bg-emerald-500"}
                                    onClick={() => {
                                        setIsModalOpen(true);
                                        setEditLocation({...resource});
                                    }}>
                                编辑
                            </button>
                            <button className={"rounded-md text-white px-2 bg-red-500"}>删除</button>
                        </td>
                    </tr>
                );
            });
            setLocationWeb(webData);
        });
    }, [locationInfo]);

    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-expect-error
    function handleInputChange(event) {
        const {name, value, type, checked} = event.target;
        setEditLocation(prevState => ({
            ...prevState,
            [name]: type === "checkbox" ? checked : value
        }));
    }

    async function editSubmit() {
        // 提交数据进行修改操作
        const getRes = await AdminEditLocationAPI(editLocation);
        if (getRes?.output === "Success") {
            message.success("操作成功");
            // 关闭模态框
            setIsModalOpen(false);
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

    return (
        <div className={"grid grid-cols-12 gap-6"}>
            <div className={"col-span-full text-xl font-bold"}>位置管理</div>
            <div className={"col-span-full lg:hidden"}>
                <MenuManagement/>
            </div>
            <div className={"col-span-8"}>
                <div className="overflow-x-auto rounded-lg shadow-md">
                    <table className="min-w-full divide-y-2 divide-gray-200 bg-white text-sm">
                        <thead className={"text-left font"}>
                        <tr>
                            <th className="whitespace-nowrap px-4 py-2 font-medium text-gray-900">名字</th>
                            <th className="whitespace-nowrap px-4 py-2 font-medium text-gray-900">别名</th>
                            <th className="whitespace-nowrap px-4 py-2 font-medium text-gray-900">是否展示</th>
                            <th className="whitespace-nowrap px-4 py-2 font-medium text-gray-900">修改时间</th>
                            <th className="whitespace-nowrap px-4 py-2 font-medium text-gray-900 text-end">操作</th>
                        </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-200">
                        {locationWeb}
                        </tbody>
                    </table>
                </div>
            </div>
            <div className={"col-span-4 hidden lg:block"}>
                <MenuManagement/>
            </div>
            <Modal title="编辑位置" open={isModalOpen} onCancel={() => setIsModalOpen(false)}
                   footer={
                       <div className={"flex gap-1 lg:gap-3 justify-end"}>
                           <button
                               className={"rounded-lg transition px-4 py-1 text-white bg-emerald-500 hover:bg-emerald-600"}
                               onClick={() => {
                                   setIsModalOpen(false);
                                   setEditLocation({} as LocationDO);
                               }}>
                               取消
                           </button>
                           <button
                               className={"rounded-lg transition px-4 py-1 text-white bg-red-400 hover:bg-red-500"}
                               onClick={editSubmit}>
                               修改
                           </button>
                       </div>
                   }>
                <div className={"grid grid-cols-2 gap-3 py-3"}>
                    <div className={"col-span-full sm:col-span-1"}>
                        <label htmlFor="LocationID" className="block text-xs font-medium text-gray-700">序列号</label>
                        <input
                            type="number"
                            id="LocationID"
                            name="id"
                            value={editLocation.id || ""}
                            className="mt-1 w-full rounded-md border-gray-200 shadow-sm text-sm bg-gray-100"
                            readOnly={true}
                        />
                    </div>
                    <div className={"col-span-full sm:col-span-1"}>
                        <label htmlFor="CreatedAt" className="block text-xs font-medium text-gray-700">创建时间</label>
                        <input
                            type="text"
                            id="CreatedAt"
                            name="created_at"
                            value={editLocation.created_at || ""}
                            className="mt-1 w-full rounded-md border-gray-200 shadow-sm text-sm bg-gray-100"
                            readOnly={true}
                        />
                    </div>
                    <div className={"col-span-full sm:col-span-1"}>
                        <label htmlFor="LocationName"
                               className="block text-xs font-medium text-gray-700">位置识别名称</label>
                        <input
                            type="text"
                            id="LocationName"
                            name="name"
                            value={editLocation.name || ""}
                            onChange={handleInputChange}
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
                            value={editLocation.display_name || ""}
                            onChange={handleInputChange}
                            className="mt-1 w-full rounded-md border-gray-200 shadow-sm text-sm"
                        />
                    </div>
                    <div className={"col-span-full"}>
                        <label htmlFor="LocationSort" className="block text-xs font-medium text-gray-700">优先级</label>
                        <input
                            type="number"
                            id="LocationSort"
                            name="sort"
                            value={editLocation.sort || ""}
                            onChange={handleInputChange}
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
                            value={editLocation.description || ""}
                            onChange={handleInputChange}
                            className="mt-2 w-full rounded-lg border-gray-200 align-top shadow-sm sm:text-sm"
                        />
                    </div>
                    <div className={"col-span-full"}>
                        <label htmlFor="Option1" className="flex cursor-pointer items-center justify-end gap-3">
                            <strong className="font-medium text-gray-900">是否对外可见</strong>
                            <input
                                type="checkbox"
                                name="reveal"
                                checked={editLocation.reveal || false}
                                onChange={handleInputChange}
                                className={"transition size-4 rounded border-gray-300 text-emerald-500"}
                            />
                        </label>
                    </div>
                </div>
            </Modal>
        </div>
    );

    function MenuManagement() {
        return (
            <div className={"rounded-lg bg-white shadow-md p-3"}>
                123
            </div>
        );
    }
}
