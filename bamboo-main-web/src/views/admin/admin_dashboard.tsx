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

import {CardBorder} from "../../components/card/card_border.tsx";
import {
    MoneyCollectOutlined,
    UserAddOutlined,
    UserDeleteOutlined,
    UserOutlined,
    UserSwitchOutlined
} from "@ant-design/icons";

export function AdminDashboard() {
    document.title = "竹叶 - 看板"

    return (
        <div className={"grid grid-cols-12 gap-6"}>
            <div className={"col-span-12 text-xl font-bold"}>首页</div>
            <div className={"col-span-12 lg:col-span-9 grid grid-cols-12 gap-6"}>
                <div className={"col-span-12 md:col-span-6"}>
                    <CardBorder type={"友链信息状态看板"} name={"友链"} element={
                        <>
                            <div className="inline-flex shrink-0 items-center gap-2">
                                <UserOutlined/>
                                <div className="mt-1.5 sm:mt-0">
                                    <p className="text-gray-500">友链总数</p>
                                    <p className="font-medium">1 个</p>
                                </div>
                            </div>
                            <div className="inline-flex shrink-0 items-center gap-2">
                                <UserAddOutlined/>
                                <div className="mt-1.5 sm:mt-0">
                                    <p className="text-gray-500">最近添加</p>
                                    <p className="font-medium">1 个</p>
                                </div>
                            </div>
                            <div className="inline-flex shrink-0 items-center gap-2">
                                <UserSwitchOutlined/>
                                <div className="mt-1.5 sm:mt-0">
                                    <p className="text-gray-500">最近修改</p>
                                    <p className="font-medium">1 个</p>
                                </div>
                            </div>
                            <div className="inline-flex shrink-0 items-center gap-2">
                                <UserDeleteOutlined/>
                                <div className="mt-1.5 sm:mt-0">
                                    <p className="text-gray-500">最近删除</p>
                                    <p className="font-medium">1 个</p>
                                </div>
                            </div>
                        </>
                    }/>
                </div>
                <div className={"col-span-12 md:col-span-6"}>
                    <CardBorder type={"赞助信息状态看板"} name={"赞助"} element={
                        <>
                            <div className="inline-flex shrink-0 items-center gap-2">
                                <MoneyCollectOutlined/>
                                <div className="mt-1.5 sm:mt-0">
                                    <p className="text-gray-500">总赞助</p>
                                    <p className="font-medium">¥ 1</p>
                                </div>
                            </div>
                            <div className="inline-flex shrink-0 items-center gap-2">
                                <UserOutlined/>
                                <div className="mt-1.5 sm:mt-0">
                                    <p className="text-gray-500">今年赞助</p>
                                    <p className="font-medium">¥ 1</p>
                                </div>
                            </div>
                            <div className="inline-flex shrink-0 items-center gap-2">
                                <UserOutlined/>
                                <div className="mt-1.5 sm:mt-0">
                                    <p className="text-gray-500">本月赞助</p>
                                    <p className="font-medium">¥ 1</p>
                                </div>
                            </div>
                            <div className="inline-flex shrink-0 items-center gap-2">
                                <UserOutlined/>
                                <div className="mt-1.5 sm:mt-0">
                                    <p className="text-gray-500">今日赞助</p>
                                    <p className="font-medium">¥ 2</p>
                                </div>
                            </div>
                        </>
                    }/>
                </div>
            </div>
            <div className={"hidden lg:block lg:col-span-3"}>
                <div
                    className="transition block bg-white rounded-lg p-4 shadow-sm shadow-indigo-100 hover:scale-105 text-center">
                    213
                </div>
            </div>
        </div>
    )
}
