// 错误码常量
export const errcode = {
  // 数据库相关错误
  DB_CONN_ERR: '1001',  // 数据库连接错误
  DB_DUP_ERR: '1062',   // 数据重复错误
  
  // JWT相关错误
  JWT_ERR: '2001',      // JWT错误
  
  // 参数相关错误
  WRONG_PARAM: '3001',  // 参数错误
  
  // 权限相关错误
  NO_PERMISSION: '4001' // 权限不足
} 