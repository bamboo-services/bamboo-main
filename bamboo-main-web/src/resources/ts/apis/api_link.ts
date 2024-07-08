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

import {LinkGetAdminEntity} from "../models/entity/link_get_admin_entity.ts";
import {BaseApi, GetAuthorizationToken, MethodType} from "../base_api.ts";
import {BaseResponse} from "../models/base_response.ts";
import {LocationDO, LocationGetAdminEntity} from "../models/entity/location_get_admin_entity.ts";
import {InnerLinkDTO, LinkGetEntity} from "../models/entity/link_get_entity.ts";
import {LocationDeleteDTO} from "../models/dto/location_delete.ts";
import {LocationAddDTO} from "../models/dto/location_add.ts";
import {ColorsEntity} from "../models/entity/color_get_entity.ts";

/**
 * GetLinkAPI
 *
 * 用于获取所有的链接，该接口只有管理员可以进行操作，其他用户不应该通过该接口进行数据的获取操作
 *
 * @returns Promise<BaseResponse<AuthLoginEntity>>
 */
const AdminGetLinkAPI = async (): Promise<BaseResponse<LinkGetAdminEntity> | undefined> => {
    return BaseApi<LinkGetAdminEntity>(
        MethodType.GET,
        "/api/v1/link/admin",
        null,
        null,
        null,
        {"Authorization": GetAuthorizationToken()}
    );
}

/**
 * GetLinkAPI
 *
 * 用于获取所有可见链接，该接口任何用户都可以进行操作获取
 *
 * @returns Promise<BaseResponse<LinkGetEntity>
 */
const GetLinkAPI = async (): Promise<BaseResponse<LinkGetEntity> | undefined> => {
    return BaseApi<LinkGetEntity>(
        MethodType.GET,
        "/api/v1/link",
        null,
        null,
        null,
        null
    );
}

/**
 * GetSingleLinkAPI
 *
 * 用于获取单个链接，该接口只允许管理员进行操作，其他用户不应该通过该接口进行数据的获取操作
 *
 * @param data 获取数据的 number 值
 * @returns Promise<BaseResponse<InnerLinkDTO>>
 */
const GetSingleLinkAPI = async (data: number): Promise<BaseResponse<InnerLinkDTO> | undefined> => {
    return BaseApi<InnerLinkDTO>(
        MethodType.GET,
        "/api/v1/link/single",
        null,
        {id: data},
        null,
        {"Authorization": GetAuthorizationToken()}
    );
}

/**
 * EditLinkAPI
 *
 * 用于编辑链接信息，该接口只有管理员可以进行操作，其他用户不应该通过该接口进行数据的获取操作
 *
 * @param data InnerLinkDTO
 * @returns Promise<BaseResponse<null>>
 */
const EditLinkAPI = async (data: InnerLinkDTO): Promise<BaseResponse<null> | undefined> => {
    return BaseApi<null>(
        MethodType.PUT,
        "/api/v1/link",
        data,
        null,
        null,
        {"Authorization": GetAuthorizationToken()}
    );

}

/**
 * AdminGetLocationAPI
 *
 * 用于获取所有的位置，该接口只有管理员可以进行操作，其他用户不应该通过该接口进行数据的获取操作
 *
 * @returns Promise<BaseResponse<LocationGetAdminEntity>>
 */
const AdminGetLocationAPI = async (): Promise<BaseResponse<LocationGetAdminEntity> | undefined> => {
    return BaseApi<LocationGetAdminEntity>(
        MethodType.GET,
        "/api/v1/link/location/full",
        null,
        null,
        null,
        {"Authorization": GetAuthorizationToken()}
    );
}

/**
 * AdminEditLocationAPI
 *
 * 用于编辑位置信息, 该接口只有管理员可以进行操作，其他用户不应该通过该接口进行数据的获取操作
 *
 * @param data LocationDO
 * @returns Promise<BaseResponse<null>>
 */
const AdminEditLocationAPI = async (data: LocationDO): Promise<BaseResponse<null> | undefined> => {
    return BaseApi<null>(
        MethodType.PUT,
        "/api/v1/link/location",
        data,
        null,
        null,
        {"Authorization": GetAuthorizationToken()}
    );
}

/**
 * AdminDelLocationAPI
 *
 * 用于删除位置信息, 该接口只有管理员可以进行操作，其他用户不应该通过该接口进行数据的获取操作
 *
 * @param data LocationDeleteDTO
 * @returns Promise<BaseResponse<null>>
 */
const AdminDelLocationAPI = async (data: LocationDeleteDTO): Promise<BaseResponse<null> | undefined> => {
    return BaseApi<null>(
        MethodType.DELETE,
        "/api/v1/link/location",
        null,
        data,
        null,
        {"Authorization": GetAuthorizationToken()}
    );
}

/**
 * AdminAddLocationAPI
 *
 * 用于添加位置信息, 该接口只有管理员可以进行操作，其他用户不应该通过该接口进行数据的获取操作
 * @param data LocationAddDTO
 * @returns Promise<BaseResponse<null>>
 */
const AdminAddLocationAPI = async (data: LocationAddDTO): Promise<BaseResponse<null> | undefined> => {
    return BaseApi<null>(
        MethodType.POST,
        "/api/v1/link/add/location",
        data,
        null,
        null,
        {"Authorization": GetAuthorizationToken()}
    );
}

/**
 * AdminGetColorAPI
 *
 * 用于获取所有的颜色，该接口只有管理员可以进行操作，其他用户不应该通过该接口进行数据的获取操作
 *
 * @returns Promise<BaseResponse<ColorsEntity>>
 */
const AdminGetColorAPI = async (): Promise<BaseResponse<ColorsEntity> | undefined> => {
    return BaseApi<ColorsEntity>(
        MethodType.GET,
        "/api/v1/link/color/full",
        null,
        null,
        null,
        {"Authorization": GetAuthorizationToken()}
    );

}

export {
    AdminGetLinkAPI,
    GetLinkAPI,
    GetSingleLinkAPI,
    EditLinkAPI,
    AdminGetLocationAPI,
    AdminEditLocationAPI,
    AdminDelLocationAPI,
    AdminAddLocationAPI,
    AdminGetColorAPI
};
